package resolve

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Constants for variable resolution.
const (
	envPrefix  = "env:"  // Prefix to identify environment variable references
	filePrefix = "file:" // Prefix to identify file references
	keyDelim   = "//"    // Delimiter to identify a key in a file
)

// ResolveVariable takes a string and resolves its value based on its prefix.
//
// If the string is prefixed with "env:", it's treated as an environment variable and resolved accordingly.
//
// If the string is prefixed with "file:", it's treated as a path to a file, optionally followed by a key
// (e.g., "file:/path/to/file//key") which specifies which line to retrieve from the file. The key is expected
// to be in the format "key = value".
//
// If no prefix is present, the string is returned as is.
//
// Parameters:
//   - value: string
//     The string to resolve.
//
// Returns:
//   - string
//     The resolved value of the input string.
//   - error
//     An error if the resolution fails.
func ResolveVariable(value string) (string, error) {
	if strings.HasPrefix(value, envPrefix) {
		return resolveEnvVariable(value[len(envPrefix):])
	}

	if strings.HasPrefix(value, filePrefix) {
		return resolveFileVariable(value[len(filePrefix):])
	}

	return value, nil
}

// resolveEnvVariable resolves a string as an environment variable name
// and retrieves its value from the environment.
//
// Parameters:
//   - envVar: string
//     The name of the environment variable.
//
// Returns:
//   - string
//     The value of the resolved environment variable.
//   - error
//     An error if the environment variable is not found.
func resolveEnvVariable(envVar string) (string, error) {
	resolvedVariable, found := os.LookupEnv(envVar)
	if !found {
		return "", fmt.Errorf("Environment variable '%s' not found.", envVar)
	}
	return resolvedVariable, nil
}

// resolveFileVariable resolves a string as a path to a file with an optional key.
// The key is expected to be in the format "key = value".
//
// Parameters:
//   - filePathWithKey: string
//     The string containing the file path and optional key.
//
// Returns:
//   - string
//     The resolved value based on the file and optional key.
//   - error
//     An error if resolving the file or key fails.
func resolveFileVariable(filePathWithKey string) (string, error) {
	lastSeparatorIndex := strings.LastIndex(filePathWithKey, keyDelim)
	filePath := filePathWithKey // default filePath (whole value)
	key := ""                   // default key (no key)

	// Check for key specification
	if lastSeparatorIndex != -1 {
		filePath = filePathWithKey[:lastSeparatorIndex]
		key = filePathWithKey[lastSeparatorIndex+len(keyDelim):]
	}

	filePath = os.ExpandEnv(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to open file '%s'. %v", filePath, err)
	}
	defer file.Close()

	if key != "" {
		return searchKeyInFile(file, key)
	}

	// No key specified, read the whole file
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Failed to read file '%s'. %v", filePath, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// searchKeyInFile searches for a specified key in a file and returns its associated value.
// The key is expected to be in the format "key = value".
//
// Parameters:
//   - file: *os.File
//     The opened file to search for the key.
//   - key: string
//     The key to search for in the file.
//
// Returns:
//   - string
//     The value associated with the specified key.
//   - error
//     An error if the key is not found in the file.
func searchKeyInFile(file *os.File, key string) (string, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		pair := strings.SplitN(line, "=", 2)
		if len(pair) == 2 && strings.TrimSpace(pair[0]) == key {
			return strings.TrimSpace(pair[1]), nil
		}
	}
	return "", fmt.Errorf("Key '%s' not found in file '%s'.", key, file.Name())
}
