package utils

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}

// IsInList checks if a value is in a list
func IsInList(value string, list []string) bool {
	for _, v := range list {
		if value == v {
			return true
		}
	}
	return false
}

// MapKeys returns the keys of a map as a slice
func MapKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// CheckFileAccessibility checks if a file exists and is accessible
func CheckFileAccessibility(filePath string) error {
	// Check if the file exists
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File does not exist: %s", filePath)
		}
		return fmt.Errorf("Error stating file '%s': %v", filePath, err)
	}

	// Try to open the file to check for readability
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Failed to open file '%s': %v", filePath, err)
	}
	file.Close() // Close immediately after opening, as we just want to check readability.

	return nil
}

// HasKey checks if a given key (or nested key) exists within a map, struct or interface.
// The function supports nested keys separated by dots, such as "key1.key2.key3".
func HasKey(s interface{}, key string) bool {
	v := reflect.ValueOf(s)
	keys := strings.Split(key, ".")

	for i, k := range keys {
		switch v.Kind() {
		case reflect.Map:
			v = v.MapIndex(reflect.ValueOf(k))
		case reflect.Struct:
			v = v.FieldByName(k)
		case reflect.Interface:
			// Extract the underlying value of the interface
			v = v.Elem()
			// Recursively call HasKey for the remaining key parts
			return HasKey(v.Interface(), strings.Join(keys[i:], "."))
		default:
			return false
		}

		if !v.IsValid() {
			return false
		}
	}
	return true
}
