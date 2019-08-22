package kubeval

import (
	"bytes"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"runtime"
	"strings"
)

func newResultError(msg string) gojsonschema.ResultError {
	r := &gojsonschema.ResultErrorFields{}
	r.SetContext(gojsonschema.NewJsonContext("error", nil))
	r.SetDescription(msg)

	return r
}

func newResultErrors(msgs []string) []gojsonschema.ResultError {
	var res []gojsonschema.ResultError
	for _, m := range msgs {
		res = append(res, newResultError(m))
	}
	return res
}

func getString(body map[string]interface{}, key string) (string, error) {
	tokens := strings.Split(key, ".")
	var obj interface{}
	var found bool
	for _, subkey := range tokens {
		obj, found = body[subkey]
		if !found {
			return "", fmt.Errorf("Missing '%s' key", key)
		}
		if obj == nil {
			return "", fmt.Errorf("Missing '%s' value", key)
		}
		body, _ = obj.(map[string]interface{})
	}
	typedValue, ok := obj.(string)
	if !ok {
		return "", fmt.Errorf("Expected string value for key '%s'", key)
	}
	return typedValue, nil
}

// detectLineBreak returns the relevant platform specific line ending
func detectLineBreak(haystack []byte) string {
	windowsLineEnding := bytes.Contains(haystack, []byte("\r\n"))
	if windowsLineEnding && runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// in is a method which tests whether the `key` is in the set
func in(set []string, key string) bool {
	for _, k := range set {
		if k == key {
			return true
		}
	}
	return false
}
