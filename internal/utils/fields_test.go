package utils

import (
	"reflect"
	"testing"
)

func TestHasFieldByPath(t *testing.T) {
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
		{"Doesn't have field in struct", TestStruct{Field1: "value1", Field2: 1}, "Field3", false},
		{"Has field in struct", TestStruct{Field1: "value1", Field2: 1}, "Field1", true},
		{"Has key in map", testMap, "key1", true},
		{"Doesn't have key in map", testMap, "key3", false},
		{"Has nested key in map", nestedMap, "level1.level2", true},
		{"Has partial nested key in map", nestedMap, "level1", true},
		{"Doesn't have nested key in map", nestedMap, "level1.level3", false},
		{"Has invalid type", []int{1, 2, 3}, "0", false}, // This should go into the default case.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := HasFieldByPath(tc.obj, tc.key)

			if got != tc.want {
				t.Fatalf("Expected %v, but got %v", tc.want, got)
			}
		})
	}
}

func TestUpdateFieldByPath(t *testing.T) {
	// Define a sample struct for testing
	type Person struct {
		Name    string
		Age     int
		Address struct {
			Street string
			City   string
		}
	}

	// Initialize a sample struct
	p := &Person{
		Name: "Alice",
		Age:  30,
		Address: struct {
			Street string
			City   string
		}{
			Street: "123 Main St",
			City:   "New York",
		},
	}

	// Test updating fields using the function
	tests := []struct {
		path      string
		newValue  interface{}
		expectErr bool
	}{
		{path: "Name", newValue: "Bob", expectErr: false},
		{path: "Age", newValue: 35, expectErr: false}, // Pass an integer for Age
		{path: "Address.Street", newValue: "456 Elm St", expectErr: false},
		{path: "Address.City", newValue: "Los Angeles", expectErr: false},
		{path: "InvalidField", newValue: "Value", expectErr: true},
		{path: "Address.InvalidField", newValue: "Value", expectErr: true},
	}

	for _, test := range tests {
		err := UpdateFieldByPath(p, test.path, test.newValue)
		if (err != nil) != test.expectErr {
			t.Errorf("UpdateFieldByPath(%s, %s) error = %v, expectErr = %v", test.path, test.newValue, err, test.expectErr)
		}
	}

	// Verify the updated struct
	expected := &Person{
		Name: "Bob",
		Age:  35,
		Address: struct {
			Street string
			City   string
		}{
			Street: "456 Elm St",
			City:   "Los Angeles",
		},
	}

	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Updated struct does not match the expected result")
	}
}

func TestGetFieldValueByPath(t *testing.T) {
	// Define a sample struct for testing
	type Person struct {
		Name    string
		Age     int
		Address struct {
			Street string
			City   string
		}
	}

	// Initialize a sample struct
	p := &Person{
		Name: "Alice",
		Age:  30,
		Address: struct {
			Street string
			City   string
		}{
			Street: "123 Main St",
			City:   "New York",
		},
	}

	// Test getting fields using the function
	tests := []struct {
		path      string
		expectVal interface{}
		found     bool
	}{
		{path: "Name", expectVal: "Alice", found: true},
		{path: "Age", expectVal: 30, found: true},
		{path: "Address.Street", expectVal: "123 Main St", found: true},
		{path: "Address.City", expectVal: "New York", found: true},
		{path: "InvalidField", expectVal: nil, found: false},
		{path: "Address.InvalidField", expectVal: nil, found: false},
	}

	for _, test := range tests {
		val, found := GetFieldValueByPath(p, test.path)
		if found != test.found {
			t.Errorf("GetFieldValueByPath(%s) found = %v, expectFound = %v", test.path, found, test.found)
		}
		if val != test.expectVal {
			t.Errorf("GetFieldValueByPath(%s) val = %v, expectVal = %v", test.path, val, test.expectVal)
		}
	}
}
