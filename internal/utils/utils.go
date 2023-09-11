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

// HasFieldByPath is a utility function designed to check if a given key exists
// within a map, struct, or interface. This function also supports checking
// for nested keys, separated by dots (e.g., "key1.key2.key3").
// Attention. Keys are case-sensitive!
func HasFieldByPath(s interface{}, path string) bool {
	keys := strings.Split(path, ".")

	return hasFieldRecursive(reflect.ValueOf(s), keys)
}

// hasFieldRecursive recursively traverses the struct and checks if the field exists.
// Returns true if the field exists, false otherwise.
func hasFieldRecursive(v reflect.Value, keys []string) bool {
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // Dereference the pointer
	}

	if len(keys) == 0 { // No more keys to check
		return true
	}

	switch v.Kind() {
	case reflect.Map:
		mapKey := reflect.ValueOf(keys[0])
		if v.MapIndex(mapKey).IsValid() {
			return hasFieldRecursive(v.MapIndex(mapKey), keys[1:])
		}
	case reflect.Struct:
		field := v.FieldByName(keys[0])
		if field.IsValid() {
			return hasFieldRecursive(field, keys[1:])
		}
	case reflect.Interface:
		return hasFieldRecursive(v.Elem(), keys)
	}

	return false
}

// UpdateFieldByPath updates a field in a struct identified by a path with a new value.
// The path is a dot-separated string that represents the hierarchy of struct fields.
// Returns an error if the field is not found or if the update fails.
func UpdateFieldByPath(data interface{}, path string, newValue interface{}) error {
	fieldNames := strings.Split(path, ".")

	return updateFieldRecursive(reflect.ValueOf(data).Elem(), fieldNames, newValue)
}

// updateFieldRecursive recursively traverses the struct and updates the field.
// Returns an error if the field is not found or if the update fails.
func updateFieldRecursive(value reflect.Value, fieldNames []string, newValue interface{}) error {
	if len(fieldNames) == 1 { // We have reached the final field
		return setField(value, fieldNames[0], newValue)
	}

	// Obtain the field
	field := value.FieldByName(fieldNames[0])
	if !field.IsValid() {
		return fmt.Errorf("No such field: %s in obj", fieldNames[0])
	}

	// Check if the field is a struct
	fieldType := field.Type()
	if fieldType.Kind() != reflect.Struct {
		return fmt.Errorf("Field %s is not a struct", fieldNames[0])
	}

	// Recursively update the nested field
	return updateFieldRecursive(field, fieldNames[1:], newValue)
}

// setField sets the value of a field in a struct.
// Returns an error if the field is not found or if the update fails.
func setField(value reflect.Value, fieldName string, newValue interface{}) error {
	field := value.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("No such field: %s in obj", fieldName)
	}

	val := reflect.ValueOf(newValue)
	valType := val.Type()

	if !valType.AssignableTo(field.Type()) {
		return fmt.Errorf("Provided value type %s cannot be assigned to field type %s", valType, field.Type())
	}

	field.Set(val)

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

// DeepCopy deep copies the source to the destination.
// The dest argument should be a pointer to the type you want to copy into.
func DeepCopy(src, dest interface{}) error {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	// Check if dest is a pointer
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return fmt.Errorf("destination is not a valid pointer")
	}

	// Create a map to keep track of copied objects to handle circular references
	copiedObjects := make(map[uintptr]reflect.Value)

	// Perform a deep copy
	deepCopy(srcValue, destValue.Elem(), copiedObjects)

	return nil
}

// deepCopy performs a deep copy of the source to the destination.
func deepCopy(src, dest reflect.Value, copiedObjects map[uintptr]reflect.Value) {
	// Check for circular references
	if src.Kind() == reflect.Ptr && copiedObjects[src.Pointer()].IsValid() {
		dest.Set(copiedObjects[src.Pointer()])
		return
	}

	if dest.Kind() != reflect.Ptr { // Destination is not a pointer
		switch src.Kind() {
		case reflect.Struct:
			// Create a new struct and add it to the copied objects map
			for i := 0; i < src.NumField(); i++ {
				srcField := src.Field(i)
				destField := dest.Field(i)
				if destField.CanSet() {
					deepCopy(srcField, destField, copiedObjects)
				}
			}
		default:
			dest.Set(src)
		}
		return
	}

	if src.Kind() != reflect.Ptr { // Source is not a pointer, create a new pointer and copy the value
		newSrc := reflect.New(src.Type())           // Create a new pointer
		deepCopy(src, newSrc.Elem(), copiedObjects) // Copy the value
		dest.Set(newSrc)                            // Set the destination to the new pointer

		return
	}

	// Both source and destination are pointers
	if src.IsNil() { // If source is nil, set destination to nil
		dest.Set(reflect.Zero(dest.Type()))
		return
	}

	// If source is not nil, create a new pointer and copy the value
	newSrc := reflect.New(src.Elem().Type())           // Create a new pointer
	copiedObjects[src.Pointer()] = newSrc              // Add the new pointer to the copied objects map
	deepCopy(src.Elem(), newSrc.Elem(), copiedObjects) // Copy the value
	dest.Set(newSrc)                                   // Set the destination to the new pointer
}
