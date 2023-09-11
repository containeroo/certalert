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

// UpdateFieldByPath updates a field in a struct identified by a path with a new value.
// The path is a dot-separated string that represents the hierarchy of struct fields.
// Returns an error if the field is not found or if the update fails.
func UpdateFieldByPath(data interface{}, path string, newValue interface{}) error {
	v := reflect.ValueOf(data)
	fieldNames := strings.Split(path, ".")

	return updateFieldRecursive(v, fieldNames, newValue)
}

// updateFieldRecursive recursively traverses the struct and updates the field's value.
// Returns an error if the field is not found or if the update fails.
func updateFieldRecursive(value reflect.Value, fieldNames []string, newValue interface{}) error {
	if value.Kind() == reflect.Ptr {
		value = value.Elem() // Dereference the pointer
	}

	if len(fieldNames) == 0 { // No more fields to update
		val := reflect.ValueOf(newValue)
		valType := val.Type()

		if !valType.AssignableTo(value.Type()) {
			return fmt.Errorf("Provided value type %s cannot be assigned to field type %s", valType, value.Type())
		}

		value.Set(val)
		return nil
	}

	switch value.Kind() {
	case reflect.Map:
		mapKey := reflect.ValueOf(fieldNames[0])
		mapValue := value.MapIndex(mapKey)
		if !mapValue.IsValid() {
			return fmt.Errorf("No such key: %s in map", fieldNames[0])
		}
		if err := updateFieldRecursive(mapValue, fieldNames[1:], newValue); err != nil {
			return err
		}
	case reflect.Interface:
		interfaceValue := value.Elem()
		if !interfaceValue.IsValid() {
			return fmt.Errorf("Nil interface encountered while traversing path")
		}
		if err := updateFieldRecursive(interfaceValue, fieldNames, newValue); err != nil {
			return err
		}
	case reflect.Struct:
		// Obtain the field
		field := value.FieldByName(fieldNames[0])
		if !field.IsValid() {
			return fmt.Errorf("No such field: %s in obj", fieldNames[0])
		}

		// Recursively update the nested field
		if err := updateFieldRecursive(field, fieldNames[1:], newValue); err != nil {
			return err
		}
	}

	return nil
}
