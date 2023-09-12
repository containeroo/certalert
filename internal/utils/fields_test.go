package utils

import (
	"reflect"
	"testing"
)

func TestHasFieldByPath(t *testing.T) {
	type Person struct {
		Name    string
		Age     int
		Address struct {
			Street string
			City   string
		}
		Items []struct {
			Name     string
			Password string
		}
	}
	p := Person{
		Name: "Alice",
		Age:  30,
		Address: struct {
			Street string
			City   string
		}{
			Street: "123 Main St",
			City:   "New York",
		},
		Items: []struct {
			Name     string
			Password string
		}{
			{"Item1", "Pass1"},
			{"Item2", "Pass2"},
		},
	}

	tests := []struct {
		path     string
		expected bool
	}{
		{"Name", true},
		{"Age", true},
		{"Address.Street", true},
		{"Items[]", true},
		{"Items[].Field1", false},
		{"Items[].Name", true},
	}

	for _, test := range tests {
		actual := HasFieldByPath(p, test.path)
		if actual != test.expected {
			t.Errorf("HasFieldByPath(%s) = %v; want %v", test.path, actual, test.expected)
		}
	}
}

func TestUpdateFieldRecursive(t *testing.T) {
	type Person struct {
		Name    string
		Age     int
		Address struct {
			Street string
			City   string
		}
		Items []struct {
			Name     string
			Password string
		}
		MoreItems []string
	}

	// Create a sample Person struct
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
		Items: []struct {
			Name     string
			Password string
		}{
			{"Item1", "Pass1"},
			{"Item2", "Pass2"},
		},
		MoreItems: []string{"Item1", "Item2"},
	}

	// Update the Password field of all items in the Items slice
	err := updateFieldRecursive(reflect.ValueOf(p), []string{"Items[]", "Password"}, "NewPassword")
	if err != nil {
		t.Fatalf("Error updating field: %v", err)
	}

	// Check if the Password field was updated for all items
	for _, item := range p.Items {
		if item.Password != "NewPassword" {
			t.Errorf("Item Password not updated correctly, expected 'NewPassword', got '%s'", item.Password)
		}
	}
}

func TestGetFieldValueByPath(t *testing.T) {
	// Define a sample struct for testing
	type Person struct {
		Name    string
		Age     int
		Empty   string
		Address struct {
			Street string
			City   string
		}
	}

	// Initialize a sample struct
	p := &Person{
		Name:  "Alice",
		Age:   30,
		Empty: "",
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
		{path: "Empty", expectVal: "", found: true},
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

func TestUpdateFieldByPathFunc(t *testing.T) {
	type Person struct {
		Name    string
		Age     int
		Address struct {
			Street string
			City   string
		}
		Items []struct {
			Name     string
			Password string
		}
	}

	// Create a sample Person struct
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
		Items: []struct {
			Name     string
			Password string
		}{
			{"Item1", "Pass1"},
			{"Item2", "Pass2"},
		},
	}

	// Update the Password field of all items in the Items slice
	err := UpdateFieldByPath(p, "Items[].Password", "NewPassword")
	if err != nil {
		t.Fatalf("Error updating field: %v", err)
	}

	// Check if the Password field was updated for all items
	for _, item := range p.Items {
		if item.Password != "NewPassword" {
			t.Errorf("Item Password not updated correctly, expected 'NewPassword', got '%s'", item.Password)
		}
	}

	// Update the Age field of the Person struct using a function
	ageIncrement := func(currentAge int) int {
		return currentAge + 1
	}

	err = UpdateFieldByPath(p, "Age", ageIncrement)
	if err != nil {
		t.Fatalf("Error updating field: %v", err)
	}

	// Check if the Age field was updated correctly
	if p.Age != 31 {
		t.Errorf("Age not updated correctly, expected 31, got %d", p.Age)
	}
}
