package utils

import (
	"bufio"
	"fmt"
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
		filePathWithKey := value[5:]
		lastSeparatorIndex := strings.LastIndex(filePathWithKey, ":")
		if lastSeparatorIndex == -1 || lastSeparatorIndex+1 == len(filePathWithKey) {
			return "", fmt.Errorf("Invalid format for file: expected 'file:path/to/file:{key}', got '%s'", value)
		}

		filePath := filePathWithKey[:lastSeparatorIndex]
		key := strings.Trim(filePathWithKey[lastSeparatorIndex+1:], "{}")

		file, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("Failed to open file '%s': %v", filePath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// Split each line into key-value pairs
			pair := strings.SplitN(line, "=", 2)
			if len(pair) == 2 && strings.TrimSpace(pair[0]) == key {
				// If the key matches the requested key, return the associated value
				return strings.TrimSpace(pair[1]), nil
			}
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("Failed to read file '%s': %v", filePath, err)
		}

		return "", fmt.Errorf("Key '%s' not found in file", key)
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
