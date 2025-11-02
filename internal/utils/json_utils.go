package utils

import (
	"encoding/json"
	"errors"
)

// ValidateJSON checks if a string is valid JSON
func ValidateJSON(data []byte) error {
	var js map[string]any
	if err := json.Unmarshal(data, &js); err != nil {
		return errors.New("invalid JSON format")
	}
	return nil
}

// PrettyPrintJSON indents JSON for readability (optional)
func PrettyPrintJSON(data any) (string, error) {
	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(pretty), nil
}
