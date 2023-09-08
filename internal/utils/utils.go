package utils

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}

// DeepCopy deep copies the source to the destination.
// The dest argument should be a pointer to the type you want to copy into.
func DeepCopy(src, dest interface{}) error {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	// Check if dest is a pointer
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return fmt.Errorf("destination is not a valid pointer")
	}

	// Check if src and dest have the same type
	if srcValue.Type() != destValue.Elem().Type() {
		return fmt.Errorf("source and destination types do not match")
	}

	// Perform a deep copy of the struct fields
	deepCopyStruct(srcValue, destValue.Elem())

	return nil
}

// deepCopyStruct recursively copies fields from srcValue to destValue.
func deepCopyStruct(srcValue, destValue reflect.Value) {
	switch srcValue.Kind() {
	case reflect.Ptr:
		if srcValue.IsNil() {
			destValue.Set(reflect.Zero(destValue.Type()))
			return
		}
		srcValue = srcValue.Elem()
		destValue = destValue.Elem()
		deepCopyStruct(srcValue, destValue)

	case reflect.Struct:
		for i := 0; i < srcValue.NumField(); i++ {
			srcField := srcValue.Field(i)
			destField := destValue.Field(i)
			if destField.CanSet() {
				deepCopyStruct(srcField, destField)
			}
		}

	default:
		destValue.Set(srcValue)
	}
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
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File does not exist: %s", filePath)
		}
		return fmt.Errorf("Failed to open file '%s': %v", filePath, err)
	}
	defer file.Close()

	return nil
}

// HasKey is a utility function designed to check if a given key exists
// within a map, struct, or interface. This function also supports checking
// for nested keys, separated by dots (e.g., "key1.key2.key3").
// Attention. Keys are case-sensitive!
func HasKey(s interface{}, key string) bool {
	v := reflect.ValueOf(s) // Obtain the Value of the passed interface{}

	// Split the key string by dots to handle nested keys
	keys := strings.Split(key, ".")

	for i, k := range keys {
		if v.Kind() == reflect.Ptr {
			// If the current object is a pointer, dereference it
			v = v.Elem()
		}

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

// ExtractHostAndPort extracts the hostname and port from the listen address
func ExtractHostAndPort(address string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, err
	}

	return host, port, nil
}

// IsValidURL tests a string to determine if it is a well-structured URL.
func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
