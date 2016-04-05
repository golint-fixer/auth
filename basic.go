package auth

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// ParseAuthHeader parses the given authentication header
func ParseAuthHeader(header string) (Token, error) {
	t := Token{}
	values := strings.Fields(header)

	// If not type based authentication scheme
	if len(values) > 2 || len(values) <= 1 {
		t.Value = header
		return t, nil
	}

	// Set token fields
	t.Value = values[1]
	t.Type = strings.ToLower(values[0])

	// Handle basic auth
	if t.Type == "basic" {
		value, err := DecodeBasicAuthHeader(values[1])
		if err != nil {
			return t, err
		}
		t.Value = value
	}

	return t, nil
}

// DecodeBasicAuthHeader decodes a given string as HTTP basic auth scheme.
func DecodeBasicAuthHeader(value string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("Failed to parse header '%s', base64 failed: %s", value, err)
	}

	hash := string(decoded)
	if len(strings.SplitN(hash, ":", 2)) != 2 {
		return "", fmt.Errorf("Failed to parse header '%s', expected separator ':'", value)
	}

	return hash, nil
}
