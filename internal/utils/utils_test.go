package utils

import (
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestIsInList(t *testing.T) {
	list := []string{"one", "two", "three"}

	if !IsInList("one", list) {
		t.Fatalf("'one' should be in list")
	}

	if IsInList("four", list) {
		t.Fatalf("'four' should not be in list")
	}
}

func TestMapKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected []string
	}{
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: nil,
		},
		{
			name:     "single key-value pair",
			input:    map[string]string{"key1": "value1"},
			expected: []string{"key1"},
		},
		{
			name:     "multiple key-value pairs",
			input:    map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
			expected: []string{"key1", "key2", "key3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractMapKeys(tt.input)
			sort.Strings(got)
			sort.Strings(tt.expected)
			if got != nil && tt.expected != nil && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("%s for input %v, expected %v, but got %v", tt.name, tt.input, tt.expected, got)
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

type NestedStruct struct {
	InnerField string
}

type TestStruct struct {
	Field1 string
	Field2 NestedStruct
}

func TestHasKey(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		key    string
		expect bool
	}{
		{
			name:   "simple map key exists",
			input:  map[string]int{"key1": 1, "key2": 2},
			key:    "key1",
			expect: true,
		},
		{
			name:   "simple map key does not exist",
			input:  map[string]int{"key1": 1, "key2": 2},
			key:    "key3",
			expect: false,
		},
		{
			name:   "nested map key exists",
			input:  map[string]map[string]int{"key1": {"nestedKey": 1}},
			key:    "key1.nestedKey",
			expect: true,
		},
		{
			name:   "nested map key does not exist",
			input:  map[string]map[string]int{"key1": {"nestedKey": 1}},
			key:    "key1.wrongKey",
			expect: false,
		},
		{
			name:   "struct key exists",
			input:  TestStruct{Field1: "value1", Field2: NestedStruct{InnerField: "inner"}},
			key:    "Field1",
			expect: true,
		},
		{
			name:   "struct key does not exist",
			input:  TestStruct{Field1: "value1", Field2: NestedStruct{InnerField: "inner"}},
			key:    "Field3",
			expect: false,
		},
		{
			name:   "nested struct key exists",
			input:  TestStruct{Field1: "value1", Field2: NestedStruct{InnerField: "inner"}},
			key:    "Field2.InnerField",
			expect: true,
		},
		{
			name:   "nested struct key does not exist",
			input:  TestStruct{Field1: "value1", Field2: NestedStruct{InnerField: "inner"}},
			key:    "Field2.WrongField",
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasKey(tt.input, tt.key); got != tt.expect {
				t.Errorf("Expected %v for key '%s', but got %v", tt.expect, tt.key, got)
			}
		})
	}
}
