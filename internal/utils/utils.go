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

// ExtractMapKeys is a utility function that takes an interface{} argument,
// checks if it's a map, and then returns the keys of that map as a slice of strings.
// If the argument is not a map or the map's keys are not strings, it returns nil.
func ExtractMapKeys(m interface{}) []string {
	v := reflect.ValueOf(m) // Get the value of m

	// Check if the value is of type 'Map'
	if v.Kind() != reflect.Map {
		return nil
	}
	keys := v.MapKeys() // Retrieve the keys of the map

	// Initialize a slice of strings to hold the keys
	strkeys := make([]string, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		// Convert the key to a string using type assertion
		keyStr, ok := keys[i].Interface().(string)
		if !ok {
			return nil
		}
		strkeys = append(strkeys, keyStr)
	}

	return strkeys
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
