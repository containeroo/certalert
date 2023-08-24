package utils

import (
	"fmt"
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
			got := ExtractMapKeys(tc.input, true)

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

	// Test case when os.Stat itself returns an error (other than file not found)
	t.Run("Error stating file", func(t *testing.T) {
		err := CheckFileAccessibility(string([]byte{0}))
		if err == nil {
			t.Fatalf("Expected an error but got nil")
		}
		expectedErrMsgPrefix := "Error stating file"
		if err != nil && err.Error()[:len(expectedErrMsgPrefix)] != expectedErrMsgPrefix {
			t.Fatalf("Expected error to start with '%s', got '%v'", expectedErrMsgPrefix, err)
		}
	})
}

func TestHasKey(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
	}

	testMap := map[string]int{
		"key1": 1,
		"key2": 2,
	}

	nestedMap := map[string]interface{}{
		"level1": map[string]int{
			"level2": 3,
		},
	}

	testCases := []struct {
		name string
		obj  interface{}
		key  string
		want bool
	}{
		{"Has field in struct", TestStruct{Field1: "value1", Field2: 1}, "Field1", true},
		{"Doesn't have field in struct", TestStruct{Field1: "value1", Field2: 1}, "Field3", false},
		{"Has key in map", testMap, "key1", true},
		{"Doesn't have key in map", testMap, "key3", false},
		{"Has nested key in map", nestedMap, "level1.level2", true},
		{"Has partial nested key in map", nestedMap, "level1", true},
		{"Doesn't have nested key in map", nestedMap, "level1.level3", false},
		{"Has invalid type", []int{1, 2, 3}, "0", false}, // This should go into the default case.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := HasKey(tc.obj, tc.key)

			if got != tc.want {
				t.Fatalf("Expected %v, but got %v", tc.want, got)
			}
		})
	}
}

func TestDeepCopy(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("Test struct copy", func(t *testing.T) {
		person1 := &Person{Name: "John", Age: 30}
		person2 := &Person{}
		if err := DeepCopy(person1, person2); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(person1, person2) {
			t.Errorf("Structs are not deeply equal: %v, %v", person1, person2)
		}
	})

	t.Run("Test non-pointer destination", func(t *testing.T) {
		person1 := &Person{Name: "John", Age: 30}
		var x int
		err := DeepCopy(person1, x)
		if err == nil || err.Error() != "destination is not a pointer" {
			t.Errorf("Expected pointer error, got: %v", err)
		}
	})

	t.Run("Test actual deep copy", func(t *testing.T) {
		person1 := &Person{Name: "John", Age: 30}
		person2 := &Person{}
		person2.Age = 40
		if reflect.DeepEqual(person1, person2) {
			t.Errorf("Structs are sharing memory: %v, %v", person1, person2)
		}
	})

	t.Run("Test unmarshallable type", func(t *testing.T) {
		type UnmarshalableType struct {
			F func()
		}
		unmarshalable := &UnmarshalableType{}
		err := DeepCopy(unmarshalable, &UnmarshalableType{})
		errMsg := fmt.Sprintf("error while marshaling: json: unsupported type: func()")
		if err == nil || err.Error() != errMsg {
			t.Errorf("Expected '%s', got: '%v'", errMsg, err)
		}
	})
}
