package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {
	// Define your config structure, for example:
	type Config struct {
		Key    string
		Nested struct {
			Key string
		}
	}

	// Temporarily create a config file
	tempFile, err := os.CreateTemp(os.TempDir(), "*.yaml")
	assert.NoError(t, err, "Failed to create temp file.")

	defer os.Remove(tempFile.Name()) // clean up

	tempFileName := tempFile.Name()

	content := []byte("key: value\nnested:\n  key: nested value\n")
	_, err = tempFile.Write(content)
	assert.NoError(t, err, "Failed to write to temp file.")
	tempFile.Close()

	var cfg Config

	// Call the function under test
	err = ReadConfigFile(tempFileName, &cfg)
	assert.NoError(t, err, "Failed to read config file.")

	// Check the values in the returned config
	assert.Equal(t, "value", cfg.Key)
	assert.Equal(t, "nested value", cfg.Nested.Key)

	// Test errors
	err = ReadConfigFile("not existing", &cfg)
	assert.Error(t, err, "Failed to read config file.")

	// marshal error
	err = ReadConfigFile(tempFileName, &cfg.Key)
	assert.Error(t, err, "Failed to read config file.")

}
