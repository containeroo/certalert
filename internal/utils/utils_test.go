package utils

import (
	"os"
	"reflect"
	"sort"
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
	list := []string{"one", "two", "three"}

	if !IsInList("one", list) {
		t.Fatalf("'one' should be in list")
	}

	if IsInList("four", list) {
		t.Fatalf("'four' should not be in list")
	}
}

func TestExtractMapKeys(t *testing.T) {
	testCases := []struct {
		name  string
		input interface{}
		want  []string
	}{
		{
			name: "Valid map with string keys",
			input: map[string]int{
				"key1": 1,
				"key2": 2,
				"key3": 3,
			},
			want: []string{"key1", "key2", "key3"},
		},
		{
			name:  "Invalid input (slice)",
			input: []int{1, 2, 3},
			want:  nil,
		},
		{
			name:  "Invalid input (string)",
			input: "hello",
			want:  nil,
		},
		{
			name: "Map with non-string keys (should fail type assertion)",
			input: map[int]string{
				1: "one",
				2: "two",
				3: "three",
			},
			want: nil, // Because your function assumes keys are strings, this will fail type assertion.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtractMapKeys(tc.input)

			sort.Strings(got)
			sort.Strings(tc.want)

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Expected %v, but got %v", tc.want, got)
			}
		})
	}
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
	tests := []struct {
		input        string
		expectedHost string
		expectedPort int
		expectedErr  bool
	}{
		{"example.com:8080", "example.com", 8080, false},
		{":1234", "", 1234, false},
		{"localhost:", "", 0, true},
		{"localhost:8080", "localhost", 8080, false},
		{"127.0.0.1:", "", 0, true},
		{"127.0.0.1:8080", "127.0.0.1", 8080, false},
		{"invalid", "", 0, true},
		{"invalid:", "", 0, true},
	}

	for _, test := range tests {
		host, port, err := ExtractHostAndPort(test.input)

		if (err != nil) != test.expectedErr {
			t.Errorf("For %s, expected error: %v, but got: %v", test.input, test.expectedErr, err != nil)
			continue
		}

		if host != test.expectedHost {
			t.Errorf("For %s, expected host: %s, but got: %s", test.input, test.expectedHost, host)
		}

		if port != test.expectedPort {
			t.Errorf("For %s, expected port: %d, but got: %d", test.input, test.expectedPort, port)
		}
	}
}

func TestIsValidURL(t *testing.T) {
	testCases := []struct {
		urlStr       string
		expectedBool bool
	}{
		{"http:/www.google.com", false}, // Malformed
		{"https://www.google.com", true},
		{"http://www.google.com", true},
		{"ftp://files.com", true},
		{"www.google.com", false}, // Missing scheme (like http, https)
		{"http://", false},        // Malformed
		{"http://10.0.0.69", true},
	}

	for _, tc := range testCases {
		t.Run(tc.urlStr, func(t *testing.T) {
			if IsValidURL(tc.urlStr) != tc.expectedBool {
				t.Errorf("Expected %v for %s", tc.expectedBool, tc.urlStr)
			}
		})
	}
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
}

func TestStructHasField(t *testing.T) {
	type InnerStruct struct {
		InnerField int
	}

	type MyStruct struct {
		Name  string
		Value int
		Inner InnerStruct
	}

	tests := []struct {
		name   string
		s      interface{}
		key    string
		expect bool
	}{
		{
			name:   "FieldNestedExists",
			s:      MyStruct{},
			key:    "Inner.InnerField",
			expect: true,
		},
		{
			name:   "FieldExists",
			s:      MyStruct{},
			key:    "Name",
			expect: true,
		},
		{
			name:   "FieldNotExists",
			s:      MyStruct{},
			key:    "NonExistentField",
			expect: false,
		},
		{
			name:   "FieldInPtr",
			s:      &MyStruct{},
			key:    "Name",
			expect: true,
		},
		{
			name:   "FieldNestedInPtr",
			s:      &MyStruct{},
			key:    "Inner.InnerField",
			expect: true,
		},
		{
			name:   "FieldInInterface",
			s:      interface{}(&MyStruct{}),
			key:    "Name",
			expect: true,
		},
		{
			name:   "FieldNestedInInterface",
			s:      interface{}(&MyStruct{}),
			key:    "Inner.InnerField",
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := HasStructField(tt.s, tt.key)
			if actual != tt.expect {
				t.Errorf("Expected %t, but got %t for key '%s'", tt.expect, actual, tt.key)
			}
		})
	}
}
