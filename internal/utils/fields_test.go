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
		{"Items[].Name", true},
		{"Items[-3].Name", false},
		{"Items[3].Name", false},
		{"Items[-1].Name", true},
		{"Items[-2].Name", true},
		{"Items[1].Name", true},
		{"Items[1].Field1", false},
		{"Name", true},
		{"Age", true},
		{"Address.Street", true},
		{"Items[]", true},
		{"Items[].Field1", false},
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
