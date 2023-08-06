package utils

import (
	"os"
	"strings"
	"testing"
)

func TestResolveVariable(t *testing.T) {
	t.Run("with env var", func(t *testing.T) {
		err := os.Setenv("TEST_VAR", "test-value")
		if err != nil {
			t.Fatalf("Failed to set environment variable: %v", err)
		}

		result, err := ResolveVariable("env:TEST_VAR")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "test-value" {
			t.Fatalf("Expected 'test-value', got '%s'", result)
		}
	})

	t.Run("with file", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "example")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write([]byte("file-test-value")); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}

		result, err := ResolveVariable("file:" + tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "file-test-value" {
			t.Fatalf("Expected 'file-test-value', got '%s'", result)
		}
	})

	t.Run("with file and key", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "example")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write([]byte("key1 = value 1\nkey2=value 2\nkey3 =   value 3")); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}

		result, err := ResolveVariable("file:" + tmpfile.Name() + ":{key2}")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "value 2" {
			t.Fatalf("Expected 'value 2', got '%s'", result)
		}
	})

	t.Run("without prefix", func(t *testing.T) {
		result, err := ResolveVariable("no-prefix-value")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "no-prefix-value" {
			t.Fatalf("Expected 'no-prefix-value', got '%s'", result)
		}
	})

	t.Run("with non-existent env var", func(t *testing.T) {
		_, err := ResolveVariable("env:NON_EXISTENT_VAR")
		if err == nil {
			t.Fatal("Expected error for non-existent environment variable, got nil")
		}
	})

	t.Run("with non-existent file", func(t *testing.T) {
		_, err := ResolveVariable("file:/path/to/non/existent/file")
		if err == nil {
			t.Fatal("Expected error for non-existent file, got nil")
		}
	})
}

func TestIsInList(t *testing.T) {
	list := []string{"one", "two", "three"}

	if !IsInList("one", list) {
		t.Fatalf("'one' should be in list")
	}

	if IsInList("four", list) {
		t.Fatalf("'four' should not be in list")
	}
}

func TestCheckFileAccessibility(t *testing.T) {
	// 1. File doesn't exist
	nonExistentPath := "./tmp/nonexistentfile12345"
	err := CheckFileAccessibility(nonExistentPath)
	if err == nil || !strings.HasPrefix(err.Error(), "File does not exist:") {
		t.Errorf("Expected a 'File does not exist' error, got '%v'", err)
	}

	// 2. File exists but isn't readable
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	os.Chmod(tmpFile.Name(), 0222) // Write-only permissions
	err = CheckFileAccessibility(tmpFile.Name())
	if err == nil {
		t.Errorf("Expected a 'failed to open file' error, got nil")
	}

	// 3. File exists and is readable
	os.Chmod(tmpFile.Name(), 0444) // Read-only permissions
	err = CheckFileAccessibility(tmpFile.Name())
	if err != nil {
		t.Errorf("Expected no error for readable file, got %v", err)
	}
}
