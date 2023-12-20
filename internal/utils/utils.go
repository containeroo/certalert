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

// Environment variable names
const (
	CertalertSilentEnv  = "CERTALERT_SILENT"
	CertalertVerboseEnv = "CERTALERT_VERBOSE"
)

// GetDebugAndTrace retrieves the debug and trace flags from environment variables.
//
// The function checks the "CERTALERT_VERBOSE" and "CERTALERT_SILENT" environment variables,
// parses their values, and returns the corresponding boolean flags.
//
// Returns:
//   - debug: bool
//     True if "CERTALERT_VERBOSE" is set to "true", false otherwise.
//   - trace: bool
//     True if "VIMBIN_SILENT" is set to "true", false otherwise.
//   - err: error
//     An error if there was an issue parsing the environment variables or if
//     both "VIMBIN_DEBUG" and "VIMBIN_TRACE" are set simultaneously.
func GetDebugAndTrace() (verbose bool, silent bool, err error) {
	// Check and parse CERTALERT_VERBOSE environment variable
	if verboseEnv := os.Getenv(CertalertVerboseEnv); verboseEnv != "" {
		verbose, err = strconv.ParseBool(verboseEnv)
		if err != nil {
			return false, false, fmt.Errorf("Unable to parse '%s'. %s", CertalertVerboseEnv, err)
		}
	}

	// Check and parse CERTALERT_SILENT environment variable
	if silentEnv := os.Getenv(CertalertSilentEnv); silentEnv != "" {
		silent, err = strconv.ParseBool(silentEnv)
		if err != nil {
			return false, false, fmt.Errorf("Unable to parse '%s'. %s", CertalertSilentEnv, err)
		}
	}

	// Check for mutual exclusivity of debug and trace
	if verbose && silent {
		return false, false, fmt.Errorf("'%s' and '%s' are mutually exclusive", CertalertVerboseEnv, CertalertSilentEnv)
	}

	return verbose, silent, nil
}

// BoolPtr returns a pointer to a bool.
//
// Parameters:
//   - b: bool
//     The boolean value to be converted to a pointer.
//
// Returns:
//   - *bool
//     A pointer to the input boolean value.
func BoolPtr(b bool) *bool {
	return &b
}

// IsInList checks if a given value is present in a list of strings.
//
// Parameters:
//   - value: string
//     The value to check for in the list.
//   - list: []string
//     The list of strings to search for the specified value.
//
// Returns:
//   - bool
//     True if the value is found in the list, false otherwise.
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
//
// Parameters:
//   - m: interface{}
//     The input value, which is expected to be a map.
//
// Returns:
//   - []string
//     A slice of strings representing the keys of the map. If the input is not a map
//     or if the map's keys are not strings, it returns nil.
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

// CheckFileAccessibility checks if a file exists and is accessible.
//
// Parameters:
//   - filePath: string
//     The path to the file to be checked for accessibility.
//
// Returns:
//   - error
//     An error indicating the accessibility status of the file. If the file does not
//     exist, the error message will indicate that. If there is an issue opening the file,
//     the error will contain information about the failure.
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

// ExtractHostAndPort extracts the hostname and port from the given listen address.
//
// Parameters:
//   - address: string
//     The listen address containing both hostname and port.
//
// Returns:
//   - string
//     The extracted hostname.
//   - int
//     The extracted port.
//   - error
//     An error indicating any issue encountered during extraction. If the address
//     is not in the expected format or if there is an error converting the port to an integer,
//     the error will provide details about the problem.
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

// IsValidURL tests whether the given string is a well-structured URL.
//
// Parameters:
//   - str: string
//     The string to be tested as a URL.
//
// Returns:
//   - bool
//     True if the string is a valid URL with a non-empty scheme and host; false otherwise.
func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// DeepCopy creates a deep copy of the source object and assigns it to the destination.
//
// Parameters:
//   - src: interface{}
//     The source object to be copied.
//   - dest: interface{}
//     The destination object to which the copied value is assigned.
//
// Returns:
//   - error
//     An error if there was an issue during the copying process.
func DeepCopy(src, dest interface{}) error {
	copied, err := copystructure.Copy(src)
	if err != nil {
		return err
	}
	// Set the value of dest to the copied value
	reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(copied))
	return nil
}

// HasStructField checks if a struct has a field with the specified key.
//
// Parameters:
//   - s: interface{}
//     The struct object to be checked for the presence of the field.
//   - key: string
//     The key representing the field, which may include nested fields separated by dots.
//
// Returns:
//   - bool
//     True if the field with the given key exists, false otherwise.
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
