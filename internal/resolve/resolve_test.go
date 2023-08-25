package resolve

import (
	"certalert/internal/test_helpers" // Make sure this path is correct

	"fmt"
	"os"
	"testing"
)

func TestResolveVariable(t *testing.T) {
	t.Run("with no variable", func(t *testing.T) {
		result, err := ResolveVariable("no-variable")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}
		if result != "no-variable" {
			t.Fatalf("Expected 'no-variable', got '%s'", result)
		}
	})

	t.Run("env var not found", func(t *testing.T) {
		_, err := ResolveVariable("env:NON_EXISTENT_ENV_VAR")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
		if err != nil && err.Error() != "Environment variable 'NON_EXISTENT_ENV_VAR' not found." {
			t.Fatalf("Expected error 'Environment variable 'NON_EXISTENT_ENV_VAR' not found.', got '%v'", err)
		}
	})

	t.Run("file not readable (fail to read)", func(t *testing.T) {
		_, err := ResolveVariable("file:/non/existing")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}

		tmpfilePath := "/non/existing"
		expectedErrMsg := fmt.Sprintf("Failed to open file '%s'. open %s: no such file or directory", tmpfilePath, tmpfilePath)
		if err != nil && err.Error() != expectedErrMsg {
			t.Fatalf("Expected error '%s', got '%v'", expectedErrMsg, err)
		}
	})

	t.Run("key not found in file", func(t *testing.T) {
		tmpfile, err := test_helpers.CreateTempFile("key1 = value 1\nkey2=value 2\nkey3 =   value 3")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		_, err = ResolveVariable("file:" + tmpfile.Name() + "//key4")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}

		expectedErrMsg := fmt.Sprintf("Key 'key4' not found in file '%s'.", tmpfile.Name())
		if err != nil && err.Error() != expectedErrMsg {
			t.Fatalf("Expected error '%s', got '%v'", expectedErrMsg, err)
		}

	})

	t.Run("with env variable", func(t *testing.T) {
		os.Setenv("TEST_ENV_VAR", "value1")
		defer os.Unsetenv("TEST_ENV_VAR")

		result, err := ResolveVariable("env:TEST_ENV_VAR")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "value1" {
			t.Fatalf("Expected 'value1', got '%s'", result)
		}
	})

	t.Run("with file variable", func(t *testing.T) {
		tmpfile, err := test_helpers.CreateTempFile("value1")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		result, err := ResolveVariable(fmt.Sprintf("file:%s", tmpfile.Name()))
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "value1" {
			t.Fatalf("Expected 'value1', got '%s'", result)
		}
	})

	t.Run("with file variable and key", func(t *testing.T) {
		tmpfile, err := test_helpers.CreateTempFile("key1 = value 1\nkey2=value 2\nkey3 =   value 3")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		result, err := resolveFileVariable(tmpfile.Name() + "//key2")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "value 2" {
			t.Fatalf("Expected 'value 2', got '%s'", result)
		}

	})
}

func TestResolveFileVariable(t *testing.T) {
	t.Run("with file and no key", func(t *testing.T) {
		tmpfile, err := test_helpers.CreateTempFile("content")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		result, err := resolveFileVariable(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "content" {
			t.Fatalf("Expected 'content', got '%s'", result)
		}
	})

	t.Run("with file and key", func(t *testing.T) {
		tmpfile, err := test_helpers.CreateTempFile("key1 = value 1\nkey2=value 2\nkey3 =   value 3")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		result, err := resolveFileVariable(tmpfile.Name() + "//key2")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "value 2" {
			t.Fatalf("Expected 'value 2', got '%s'", result)
		}
	})

	t.Run("with file and key with spaces and tabs", func(t *testing.T) {
		tmpfile, err := test_helpers.CreateTempFile("key1 =	  value 1\nkey2=value 2\nkey3 =   value 3")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		result, err := resolveFileVariable(tmpfile.Name() + "//key1")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "value 1" {
			t.Fatalf("Expected 'value 1', got '%s'", result)
		}
	})

	t.Run("with non-existent file", func(t *testing.T) {
		_, err := resolveFileVariable("non-existent-file")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	t.Run("with empty key", func(t *testing.T) {
		_, err := resolveFileVariable("non-existent-file//")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	t.Run("with empty key and no file", func(t *testing.T) {
		_, err := resolveFileVariable("{}")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	t.Run("with empty key and empty file", func(t *testing.T) {
		_, err := resolveFileVariable("//")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	t.Run("with colon in file path", func(t *testing.T) {
		_, err := resolveFileVariable("file:/path/to/file:with:colon")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	t.Run("with colon in key", func(t *testing.T) {
		_, err := resolveFileVariable("file:/path/to/file//key:with:colon")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	t.Run("with colon in file path and key", func(t *testing.T) {
		_, err := resolveFileVariable("file:/path/to/file:with:colon//key:with:colon")
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})
}

func TestSearchKeyInFile(t *testing.T) {

	t.Run("with file and key", func(t *testing.T) {
		tmpfile, err := test_helpers.CreateTempFile("key1 =	value 1\nkey2=value 2\nkey3 =   value 3")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		file, err := test_helpers.ReadFile(tmpfile.Name())
		defer file.Close()

		if err != nil {
			t.Fatalf("Failed to read temp file: %v", err)
		}

		result, err := searchKeyInFile(file, "key2")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "value 2" {
			t.Fatalf("Expected 'value 2', got '%s'", result)
		}
	})

	t.Run("with file and key with spaces and tabs", func(t *testing.T) {
		tmpfile, err := test_helpers.CreateTempFile("key1 =	value 1\nkey2	  =   	  value 2\nkey3 =   value 3")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		file, err := test_helpers.ReadFile(tmpfile.Name())
		defer file.Close()

		if err != nil {
			t.Fatalf("Failed to read temp file: %v", err)
		}

		result, err := searchKeyInFile(file, "key2")
		if err != nil {
			t.Fatalf("Failed to resolve variable: %v", err)
		}

		if result != "value 2" {
			t.Fatalf("Expected 'value 2', got '%s'", result)
		}
	})
}
