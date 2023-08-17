package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	envPrefix     = "env:"  // Prefix to identify environment variable references
	filePrefix    = "file:" // Prefix to identify file references
	keyDelim      = ":"     // Delimiter to identify a key in a file
	keyStartDelim = "{"     // Starting delimiter for a key in a file reference
	keyEndDelim   = "}"     // Ending delimiter for a key in a file reference
)

// ResolveVariable takes a string and resolves its value based on its prefix.
// If the string is prefixed with "env:", it's treated as an environment variable and resolved accordingly.
// If the string is prefixed with "file:", it's treated as a path to a file, optionally followed by a key
// (e.g., "file:/path/to/file:{key}") which specifies which line to retrieve from the file.
// If no prefix is present, the string is returned as is.
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
func resolveEnvVariable(envVar string) (string, error) {
	resolvedVariable, found := os.LookupEnv(envVar)
	if !found {
		return "", fmt.Errorf("Environment variable '%s' not found", envVar)
	}
	return resolvedVariable, nil
}

// resolveFileVariable resolves a string as a path to a file with an optional key.
// The key, if present, specifies which line to retrieve from the file.
func resolveFileVariable(filePathWithKey string) (string, error) {
	lastSeparatorIndex := strings.LastIndex(filePathWithKey, keyDelim)
	filePath := filePathWithKey // default filePath (whole value)
	key := ""                   // default key (no key)

	// Check if there's a key specification in the filePathWithKey
	if lastSeparatorIndex != -1 && lastSeparatorIndex+1 < len(filePathWithKey) &&
		strings.HasPrefix(filePathWithKey[lastSeparatorIndex+1:], keyStartDelim) {
		filePath = filePathWithKey[:lastSeparatorIndex]
		key = strings.Trim(filePathWithKey[lastSeparatorIndex+2:], keyEndDelim)
		if key == "" {
			return "", fmt.Errorf("Empty key")
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to open file '%s': %v", filePath, err)
	}
	defer file.Close()

	if key != "" {
		return searchKeyInFile(file, key)
	}

	// No key specified, return the whole file content
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Failed to read file '%s': %v", filePath, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// searchKeyInFile searches for a specified key in a file and returns its associated value.
// The file is expected to contain lines in the format "key=value".
func searchKeyInFile(file *os.File, key string) (string, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		pair := strings.SplitN(line, "=", 2)
		if len(pair) == 2 && strings.TrimSpace(pair[0]) == key {
			return strings.TrimSpace(pair[1]), nil
		}
	}
	return "", fmt.Errorf("Key '%s' not found in file", key)
}
