package test_helpers

import (
	"io/ioutil"
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

func TestReadFile(t *testing.T) {
	t.Run("successful file read", func(t *testing.T) {
		content := "hello, world"
		tmpfile, err := ioutil.TempFile("", "example.*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpfile.Close()

		file, err := ReadFile(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		file.Close()
	})

	t.Run("fail to read non-existent file", func(t *testing.T) {
		_, err := ReadFile("/path/to/non/existent/file")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})
}
