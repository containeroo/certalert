package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ResolveVariable resolves the value of a variable that is prefixed with env: or file:
func ResolveVariable(value string) (string, error) {
	if strings.HasPrefix(value, "env:") {
		envVar := value[4:]

		resolvedVariable, found := os.LookupEnv(envVar)
		if !found {
			return "", fmt.Errorf("Environment variable '%s' not found", envVar)
		}

		return resolvedVariable, nil
	}

	if strings.HasPrefix(value, "file:") {
		filePathWithKey := value[5:]                                  // remove the file: prefix
		lastSeparatorIndex := strings.LastIndex(filePathWithKey, ":") // find the last colon
		filePath := filePathWithKey                                   // by default, the file path is the whole value
		key := ""

		// If the value contains a colon, we split the value into a file path and a key
		if lastSeparatorIndex != -1 && lastSeparatorIndex+1 < len(filePathWithKey) {
			filePath = filePathWithKey[:lastSeparatorIndex]                  // everything before the last colon is the file path
			key = strings.Trim(filePathWithKey[lastSeparatorIndex+1:], "{}") // everything after the last colon is the key
		}

		file, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("Failed to open file '%s': %v", filePath, err)
		}
		defer file.Close()

		// If a key is specified, we look for it in the file content
		if key != "" {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pair := strings.SplitN(line, "=", 2)
				if len(pair) == 2 && strings.TrimSpace(pair[0]) == key {
					// If the key matches the requested key, return the associated value
					return strings.TrimSpace(pair[1]), nil
				}
			}
			// If we reach this point, the key was not found in the file
			return "", fmt.Errorf("Key '%s' not found in file", key)
		}

		// If no key is specified, read the whole content of the file
		data, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("Failed to read file '%s': %v", filePath, err)
		}
		return strings.TrimSpace(string(data)), nil
	}

	// if the value is not prefixed with env: or file: then just return the value
	return value, nil
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
