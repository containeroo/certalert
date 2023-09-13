package test_helpers

import (
	"os"
	"testing"
)

func TestCreateTempFile(t *testing.T) {
	t.Run("successful file creation", func(t *testing.T) {
		content := "hello, world"
		tmpfile, err := CreateTempFile(content)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		data, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to read temp file: %v", err)
		}

		if string(data) != content {
			t.Errorf("Expected file content to be '%s', got '%s'", content, string(data))
		}
	})
}
