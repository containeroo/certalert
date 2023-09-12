package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// parseInt parses an integer from a string.
func parseInt(s string) (i int) {
	fmt.Sscanf(s, "%d", &i)
	return i
}

// parseField parses a field name and returns the index and field name.
// If the field name does not contain an index, index is nil.
func parseField(name string) (index *int, fieldName string) {
	if strings.HasSuffix(name, "[]") {
		return nil, name[:len(name)-2]
	}

	if !strings.Contains(name, "[") && !strings.Contains(name, "]") {
		return nil, name
	}

	fieldName = name[:strings.Index(name, "[")] // Remove everything after the first '['
	i := parseInt(name[strings.Index(name, "[")+1 : strings.Index(name, "]")])

	index = &i
	return index, fieldName
}

// HasFieldByPath is a utility function designed to check if a given key exists
// within a map, struct, or interface. This function also supports checking
// for nested keys, separated by dots (e.g., "key1.key2.key3").
// To check if a field exists within a slice or array, append "[]" to the field name.
// Returns true if the field exists, false otherwise.
// Attention. Keys are case-sensitive!
func HasFieldByPath(data interface{}, path string) bool {
	v := reflect.ValueOf(data)
	keys := strings.Split(path, ".")

	return hasFieldRecursive(v, keys)
}

// hasFieldRecursive recursively traverses the struct and checks if the field exists.
// Returns true if the field exists, false otherwise.
func hasFieldRecursive(value reflect.Value, keys []string) bool {
	if value.Kind() == reflect.Ptr {
		value = value.Elem() // Dereference the pointer
	}

	if len(keys) == 0 { // No more keys to check
		return true
	}

	idx, fieldName := parseField(keys[0])

	switch value.Kind() {
	case reflect.Map:
		mapKey := reflect.ValueOf(fieldName)
		if value.MapIndex(mapKey).IsValid() {
			return hasFieldRecursive(value.MapIndex(mapKey), keys[1:])
		}
	case reflect.Struct:
		// Locate the field within the struct
		field := value.FieldByName(fieldName)
		if !field.IsValid() {
			return false
		}

		// Check if the field is a slice or array
		if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
			startIdx, endIdx := 0, field.Len() // Default to all elements

			if idx != nil {
				if *idx < 0 {
					*idx += endIdx // Convert to positive index. This is possible due zero-based index.
				}
				if *idx < 0 || *idx >= endIdx { // Index out of bounds
					return false
				}
				startIdx, endIdx = *idx, *idx+1 // Only check the specified index
			}

			for i := startIdx; i < endIdx; i++ {
				sliceValue := field.Index(i)
				if hasFieldRecursive(sliceValue, keys[1:]) {
					return true
				}
			}
			return false
		}
		return hasFieldRecursive(field, keys[1:])
	case reflect.Interface:
		return hasFieldRecursive(value.Elem(), keys)
	}

	return false
}

// UpdateFieldByPath updates a field in a struct identified by a path with a new value or function.
// The path is a dot-separated string that represents the hierarchy of struct fields.
// To iterate over all elements of a slice or array, append "[]" to the field name.
// If newValue is a function, it must have the signature func(T) T where T is assignable to the field type.
// Returns an error if the field is not found or if the update fails.
func UpdateFieldByPath(data interface{}, path string, newValue interface{}) error {
	value := reflect.ValueOf(data)
	fieldNames := strings.Split(path, ".")

	// Check if newValue is a function
	if reflect.TypeOf(newValue).Kind() == reflect.Func {
		// If newValue is a function, check its signature
		// The function must have the signature func(T) T where T is assignable to the field type
		newValueFunc := reflect.ValueOf(newValue)
		if newValueFunc.Type().NumIn() != 1 || // Must have one input parameter
			newValueFunc.Type().Out(0) != newValueFunc.Type().In(0) { // Must return a value of the same type as the input
			return fmt.Errorf("Function signature must be func(%s) %s", value.Type(), value.Type())
		}
	}

	return updateFieldRecursive(value, fieldNames, newValue)
}

// updateFieldRecursive recursively traverses the struct and updates the field's value.
// Returns an error if the field is not found, if the update fails, or if newValue is not a valid function.
func updateFieldRecursive(value reflect.Value, fieldNames []string, newValue interface{}) error {
	if value.Kind() == reflect.Ptr {
		value = value.Elem() // Dereference the pointer
	}

	if len(fieldNames) == 0 { // No more fields to update
		if reflect.TypeOf(newValue).Kind() == reflect.Func {
			// If newValue is a function, check its signature
			newValueFunc := reflect.ValueOf(newValue)
			if newValueFunc.Type().NumIn() != 1 || // Must have one input parameter
				newValueFunc.Type().Out(0) != value.Type() { // Must return a value of the same type as the input
				return fmt.Errorf("Function signature must be func(%s) %s", value.Type(), value.Type())
			}

			// Call the function with the current field value to get the new value
			newValueResult := newValueFunc.Call([]reflect.Value{value})[0]

			if !newValueResult.Type().AssignableTo(value.Type()) {
				return fmt.Errorf("Function result type %s cannot be assigned to field type %s", newValueResult.Type(), value.Type())
			}

			value.Set(newValueResult)
			return nil
		}

		val := reflect.ValueOf(newValue)
		valType := val.Type()

		if !valType.AssignableTo(value.Type()) {
			return fmt.Errorf("Provided value type %s cannot be assigned to field type %s", valType, value.Type())
		}

		value.Set(val)
		return nil
	}

	idx, fieldName := parseField(fieldNames[0])

	switch value.Kind() {
	case reflect.Map:
		mapKey := reflect.ValueOf(fieldName)
		mapValue := value.MapIndex(mapKey)
		if !mapValue.IsValid() {
			return fmt.Errorf("No such key: %s in map", fieldName)
		}

		return updateFieldRecursive(mapValue, fieldNames[1:], newValue)
	case reflect.Interface:
		interfaceValue := value.Elem()
		if !interfaceValue.IsValid() {
			return fmt.Errorf("Nil interface encountered while traversing path")
		}

		return updateFieldRecursive(interfaceValue, fieldNames, newValue)
	case reflect.Struct:
		// Locate the field within the struct
		field := value.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("No such field: %s in obj", fieldName)
		}

		// Check if the field is a slice or array
		if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
			startIdx, endIdx := 0, field.Len() // Default to all elements

			if idx != nil {
				if *idx < 0 {
					*idx += endIdx // Convert to positive index. This is possible due zero-based index.
				}
				if *idx < 0 || *idx >= endIdx { // Index out of bounds
					return nil
				}
				startIdx, endIdx = *idx, *idx+1 // Only check the specified index
			}

			for i := startIdx; i < endIdx; i++ {
				sliceValue := field.Index(i)
				if err := updateFieldRecursive(sliceValue, fieldNames[1:], newValue); err != nil {
					return err
				}
			}
			// If reached here, the field was updated successfully
			return nil
		}

		return updateFieldRecursive(field, fieldNames[1:], newValue)
	}

	return nil
}
