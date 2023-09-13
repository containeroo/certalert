package utils

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/copystructure"
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

// DeepCopy copies the value of src to dest
func DeepCopy(src, dest interface{}) error {
	copied, err := copystructure.Copy(src)
	if err != nil {
		return err
	}
	// Set the value of dest to the copied value
	reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(copied))

	return nil
}

// HasStructField checks if a struct has a field with the given key
func HasStructField(s interface{}, key string) bool {
	v := reflect.ValueOf(s) // Obtain the Value of the passed interface{}

	keys := strings.Split(key, ".") // Split the key string by dots to handle nested keys

	for i, k := range keys {
		if v.Kind() == reflect.Ptr { // If the current object is a pointer, dereference it
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			v = v.FieldByName(k) // Retrieve the field with the name corresponding to the key
		case reflect.Interface:
			v = v.Elem() // Dereference the interface to get its underlying value
			// Use recursion to continue checking for the remaining nested keys
			return HasStructField(v.Interface(), strings.Join(keys[i:], "."))
		default:
			return false
		}

		if !v.IsValid() {
			return false
		}
	}

	return true
}
