package utils

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestBoolPtr(t *testing.T) {
	t.Run("returns pointer to true", func(t *testing.T) {
		b := true
		result := BoolPtr(b)

		if result == nil {
			t.Fatalf("Expected a non-nil pointer, got nil")
		}

		if *result != b {
			t.Fatalf("Expected pointer to point to %v, got %v", b, *result)
		}
	})

	t.Run("returns pointer to false", func(t *testing.T) {
		b := false
		result := BoolPtr(b)

		if result == nil {
			t.Fatalf("Expected a non-nil pointer, got nil")
		}

		if *result != b {
			t.Fatalf("Expected pointer to point to %v, got %v", b, *result)
		}
	})
}

func TestIsInList(t *testing.T) {
	t.Run("element exists in the list", func(t *testing.T) {
		list := []string{"one", "two", "three"}
		element := "one"

		if !IsInList(element, list) {
			t.Fatalf("'%s' should be in the list", element)
		}
	})

	t.Run("element does not exist in the list", func(t *testing.T) {
		list := []string{"one", "two", "three"}
		element := "four"

		if IsInList(element, list) {
			t.Fatalf("'%s' should not be in the list", element)
		}
	})
}

func TestExtractMapKeys(t *testing.T) {
	t.Run("Valid map with string keys", func(t *testing.T) {
		input := map[string]int{
			"key1": 1,
			"key2": 2,
			"key3": 3,
		}
		expected := []string{"key1", "key2", "key3"}
		result := ExtractMapKeys(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
	t.Run("Valid map with int keys", func(t *testing.T) {
		input := map[int]string{
			1: "one",
			2: "two",
			3: "three",
		}
		expected := []string(nil)
		result := ExtractMapKeys(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
	t.Run("Valid map with mixed keys", func(t *testing.T) {
		input := map[interface{}]string{
			1: "one",
			2: "two",
			3: "three",
		}
		expected := []string(nil)
		result := ExtractMapKeys(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
	t.Run("Invalid input", func(t *testing.T) {
		input := "not a map"
		expected := []string(nil)
		result := ExtractMapKeys(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestCheckFileAccessibility(t *testing.T) {
	t.Run("File doesn't exist", func(t *testing.T) {
		nonExistentPath := "./tmp/nonexistentfile12345"
		err := CheckFileAccessibility(nonExistentPath)
		if err == nil || !strings.HasPrefix(err.Error(), "File does not exist:") {
			t.Errorf("Expected a 'File does not exist' error, got '%v'", err)
		}
	})

	t.Run("File exists but isn't readable", func(t *testing.T) {
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
	})

	t.Run("File exists and is readable", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		os.Chmod(tmpFile.Name(), 0444) // Read-only permissions
		err = CheckFileAccessibility(tmpFile.Name())
		if err != nil {
			t.Errorf("Expected no error for readable file, got %v", err)
		}
	})
}

func TestExtractHostAndPort(t *testing.T) {
	t.Run("Empty input", func(t *testing.T) {
		host, port, err := ExtractHostAndPort("")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err.Error() != "missing port in address" {
			t.Errorf("Expected 'Invalid host:port format', got '%s'", err.Error())
		}

		if host != "" {
			t.Errorf("Expected empty host, got %s", host)
		}
		if port != 0 {
			t.Errorf("Expected port 0, got %d", port)
		}
	})
	t.Run("Invalid input", func(t *testing.T) {
		host, port, err := ExtractHostAndPort("invalid")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err.Error() != "address invalid: missing port in address" {
			t.Errorf("Expected 'Invalid host:port format', got '%s'", err.Error())
		}

		if host != "" {
			t.Errorf("Expected empty host, got %s", host)
		}
		if port != 0 {
			t.Errorf("Expected port 0, got %d", port)
		}
	})

	t.Run("Valid input", func(t *testing.T) {
		host, port, err := ExtractHostAndPort("example.com:8080")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if host != "example.com" {
			t.Errorf("Expected host 'example.com', got %s", host)
		}
		if port != 8080 {
			t.Errorf("Expected port 8080, got %d", port)
		}
	})
}

func TestIsValidURL(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		if IsValidURL("") {
			t.Errorf("Expected false for empty string")
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		if IsValidURL("http://") {
			t.Errorf("Expected false for invalid URL")
		}
	})

	t.Run("Valid URL", func(t *testing.T) {
		if !IsValidURL("http://www.google.com") {
			t.Errorf("Expected true for valid URL")
		}
	})
	t.Run("Valid URL with port", func(t *testing.T) {
		if !IsValidURL("http://www.google.com:8080") {
			t.Errorf("Expected true for valid URL")
		}
	})

	t.Run("Valid URL with path", func(t *testing.T) {
		if !IsValidURL("http://www.google.com/path") {
			t.Errorf("Expected true for valid URL")
		}
	})

	t.Run("Valid URL with query", func(t *testing.T) {
		if !IsValidURL("http://www.google.com?query=1") {
			t.Errorf("Expected true for valid URL")
		}
	})

	t.Run("Valid URL with fragment", func(t *testing.T) {
		if !IsValidURL("http://www.google.com#fragment") {
			t.Errorf("Expected true for valid URL")
		}
	})

	t.Run("Valid URL with path, query and fragment", func(t *testing.T) {
		if !IsValidURL("http://www.google.com/path?query=1#fragment") {
			t.Errorf("Expected true for valid URL")
		}
	})
}

func TestDeepCopyNestedStruct(t *testing.T) {
	type SimpleStruct struct {
		Field1 string
		Field2 int
	}
	type NestedStruct struct {
		Field1 string
		Field2 SimpleStruct
	}

	t.Run("deep copy nested struct", func(t *testing.T) {
		src := NestedStruct{
			Field1: "Outer",
			Field2: SimpleStruct{"Hello", 42},
		}
		var dest NestedStruct

		err := DeepCopy(src, &dest)
		if err != nil {
			t.Fatalf("Error during DeepCopy: %v", err)
		}

		if !reflect.DeepEqual(src, dest) {
			t.Errorf("DeepCopy result does not match source.\nSource: %+v\nDest: %+v", src, dest)
		}

		// Change the source and make sure the dest doesn't change
		src.Field2.Field1 = "Goodbye"
		if reflect.DeepEqual(src, dest) {
			t.Errorf("DeepCopy result should not match source.\nSource: %+v\nDest: %+v", src, dest)
		}

		// Change the dest and make sure the source doesn't change
		dest.Field2.Field1 = "Also Goodbye"
		if reflect.DeepEqual(src, dest) {
			t.Errorf("DeepCopy result should not match source.\nSource: %+v\nDest: %+v", src, dest)
		}
	})
}

func TestHasStructField(t *testing.T) {
	type InnerStruct struct {
		InnerField int
	}

	type MyStruct struct {
		Name  string
		Value int
		Inner InnerStruct
	}

	t.Run("field nested exists", func(t *testing.T) {
		s := MyStruct{}
		key := "Inner.InnerField"
		expect := true

		actual := HasStructField(s, key)

		if actual != expect {
			t.Errorf("Expected %t, but got %t for key '%s'", expect, actual, key)
		}
	})

	t.Run("field exists", func(t *testing.T) {
		s := MyStruct{}
		key := "Name"
		expect := true

		actual := HasStructField(s, key)

		if actual != expect {
			t.Errorf("Expected %t, but got %t for key '%s'", expect, actual, key)
		}
	})

	t.Run("field not exists", func(t *testing.T) {
		s := MyStruct{}
		key := "NonExistentField"
		expect := false

		actual := HasStructField(s, key)

		if actual != expect {
			t.Errorf("Expected %t, but got %t for key '%s'", expect, actual, key)
		}
	})

	t.Run("field in pointer", func(t *testing.T) {
		s := &MyStruct{}
		key := "Name"
		expect := true

		actual := HasStructField(s, key)

		if actual != expect {
			t.Errorf("Expected %t, but got %t for key '%s'", expect, actual, key)
		}
	})

	t.Run("field nested in pointer", func(t *testing.T) {
		s := &MyStruct{}
		key := "Inner.InnerField"
		expect := true

		actual := HasStructField(s, key)

		if actual != expect {
			t.Errorf("Expected %t, but got %t for key '%s'", expect, actual, key)
		}
	})

	t.Run("field in interface", func(t *testing.T) {
		s := interface{}(&MyStruct{})
		key := "Name"
		expect := true

		actual := HasStructField(s, key)

		if actual != expect {
			t.Errorf("Expected %t, but got %t for key '%s'", expect, actual, key)
		}
	})

	t.Run("field nested in interface", func(t *testing.T) {
		s := interface{}(&MyStruct{})
		key := "Inner.InnerField"
		expect := true

		actual := HasStructField(s, key)

		if actual != expect {
			t.Errorf("Expected %t, but got %t for key '%s'", expect, actual, key)
		}
	})
}
