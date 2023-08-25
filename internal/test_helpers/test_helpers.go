package test_helpers

import (
	"fmt"
	"os"
)

// CreateTempFile creates a temporary file with the given content and returns a pointer to it.
func CreateTempFile(content string) (*os.File, error) {
	tmpfile, err := os.CreateTemp("", "example.*.yaml")
	if err != nil {
		return nil, fmt.Errorf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		return nil, fmt.Errorf("Failed to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		return nil, fmt.Errorf("Failed to close temp file: %v", err)
	}

	return tmpfile, nil
}

// ReadFile opens a file and returns a pointer to it.
func ReadFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file '%s': %v", filePath, err)
	}
	return file, nil
}
