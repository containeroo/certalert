package config

import (
	"certalert/internal/test_helpers"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("Valid full configuration", func(t *testing.T) {
		filePath, err := test_helpers.CreateTempFile("autoReloadConfig: true\nversion: \"1.0\"")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(filePath.Name())

		fc, _ := os.ReadFile(filePath.Name())
		fmt.Println(string(fc))

		config := &Config{}
		err = config.Read(filePath.Name())
		assert.NoError(t, err)
		assert.Equal(t, true, config.AutoReloadConfig)
		assert.Equal(t, "1.0", config.Version)
	})

	t.Run("Partial configuration", func(t *testing.T) {
		filePath, err := test_helpers.CreateTempFile("version: \"1.1\"")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(filePath.Name())

		config := &Config{}
		err = config.Read(filePath.Name())
		assert.NoError(t, err)
		assert.Equal(t, false, config.AutoReloadConfig)
		assert.Equal(t, "1.1", config.Version)
	})

	// Test reading a non-existent file
	t.Run("Invalid file path", func(t *testing.T) {
		config := &Config{}
		err := config.Read("path/to/non_existent.yaml")
		assert.Error(t, err)
	})

	t.Run("Invalid YAML", func(t *testing.T) {
		filePath, err := test_helpers.CreateTempFile("FailOnError: not_an_integer")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(filePath.Name())

		// Run the test
		cfg := &Config{}
		err = cfg.Read(filePath.Name())
		assert.Error(t, err) // We expect an error because the file has incorrect content
		assert.Contains(t, err.Error(), "Failed to unmarshal config file: 1 error(s) decoding:\n\n* cannot parse 'failOnError' as bool: strconv.ParseBool: parsing \"not_an_integer\": invalid syntax")
	})
}
