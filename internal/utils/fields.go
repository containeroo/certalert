package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// HasFieldByPath is a utility function designed to check if a given key exists
// within a map, struct, or interface. This function also supports checking
// for nested keys, separated by dots (e.g., "key1.key2.key3").
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

	// Remove "[]" suffix from field name, if present
	fieldName := keys[0]
	if strings.HasSuffix(fieldName, "[]") {
		fieldName = strings.TrimSuffix(fieldName, "[]")
	}

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
			// Iterate over the slice elements and update them
			for i := 0; i < field.Len(); i++ {
				sliceValue := field.Index(i)
				if hasFieldRecursive(sliceValue, keys[1:]) {
					return true
				}
				return false
			}
		}
		return hasFieldRecursive(field, keys[1:])
	case reflect.Interface:
		return hasFieldRecursive(value.Elem(), keys)
	}

	return false
}

// GetFieldValueByPath retrieves the value of a field in a struct identified by a path.
// The path is a dot-separated string that represents the hierarchy of struct fields.
// Returns the field's value and a boolean indicating whether the field exists or not.
func GetFieldValueByPath(data interface{}, path string) (interface{}, bool) {
	v := reflect.ValueOf(data)
	keys := strings.Split(path, ".")

	return getFieldValueRecursive(v, keys)
}

// getFieldValueRecursive recursively traverses the struct and retrieves the field's value.
// Returns the field's value and a boolean indicating whether the field exists or not.
func getFieldValueRecursive(value reflect.Value, keys []string) (interface{}, bool) {
	if value.Kind() == reflect.Ptr {
		value = value.Elem() // Dereference the pointer
	}

	if len(keys) == 0 { // No more keys to traverse
		return value.Interface(), true
	}

	switch value.Kind() {
	case reflect.Map:
		mapKey := reflect.ValueOf(keys[0])
		if value.MapIndex(mapKey).IsValid() {
			return getFieldValueRecursive(value.MapIndex(mapKey), keys[1:])
		}
	case reflect.Struct:
		field := value.FieldByName(keys[0])
		if field.IsValid() {
			return getFieldValueRecursive(field, keys[1:])
		}
	case reflect.Interface:
		return getFieldValueRecursive(value.Elem(), keys)
	}

	return nil, false
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

	// Remove "[]" suffix from field name, if present
	fieldName := fieldNames[0]
	if strings.HasSuffix(fieldName, "[]") {
		fieldName = strings.TrimSuffix(fieldName, "[]")
	}

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
			// Iterate over the slice elements and update them
			for i := 0; i < field.Len(); i++ {
				sliceValue := field.Index(i)
				if err := updateFieldRecursive(sliceValue, fieldNames[1:], newValue); err != nil {
					return err
				}
			}
			return nil
		}
		return updateFieldRecursive(field, fieldNames[1:], newValue)
	}

	return nil
}
