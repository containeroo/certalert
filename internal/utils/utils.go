package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}

// DeepCopy deep copies the source to the destination.
// The dest argument should be a pointer to the type you want to copy into.
func DeepCopy(src interface{}, dest interface{}) error {
	// Marshal the source object to a JSON byte array
	bytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("error while marshaling: %v", err)
	}

	// Make sure dest is a pointer
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return fmt.Errorf("destination is not a pointer")
	}

	// Unmarshal the JSON byte array to the destination object
	err = json.Unmarshal(bytes, dest)
	if err != nil {
		return fmt.Errorf("error while unmarshaling: %v", err)
	}

	return nil
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

// HasKey is a utility function designed to check if a given key exists
// within a map, struct, or interface. This function also supports checking
// for nested keys, separated by dots (e.g., "key1.key2.key3").
func HasKey(s interface{}, key string) bool {
	v := reflect.ValueOf(s) // Obtain the Value of the passed interface{}

	// Split the key string by dots to handle nested keys
	keys := strings.Split(key, ".")

	for i, k := range keys {
		switch v.Kind() {
		// If the current object is a map
		case reflect.Map: // Look for the key within the map
			v = v.MapIndex(reflect.ValueOf(k)) // Look for the key within the map
		case reflect.Struct: // If the current object is a struct
			v = v.FieldByName(k) // Retrieve the field with the name corresponding to the key
		case reflect.Interface: // If the current object is an interface
			v = v.Elem() // Dereference the interface to get its underlying value
			// Use recursion to continue checking for the remaining nested keys
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
