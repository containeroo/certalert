package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	// Helper to create temporary files
	createTempFile := func(content string) string {
		tempFile, err := ioutil.TempFile(os.TempDir(), "*.yaml")
		assert.NoError(t, err)
		defer tempFile.Close()

		_, err = tempFile.WriteString(content)
		assert.NoError(t, err)

		return tempFile.Name()
	}

	// Test reading a valid full configuration
	t.Run("Valid full configuration", func(t *testing.T) {
		filePath := createTempFile("autoReloadConfig: true\nversion: \"1.0\"")
		defer os.Remove(filePath)

		config := &Config{}
		err := config.Read(filePath)
		assert.NoError(t, err)
		assert.Equal(t, true, config.AutoReloadConfig)
		assert.Equal(t, "1.0", config.Version)
	})

	// Test reading a partial configuration
	t.Run("Partial configuration", func(t *testing.T) {
		filePath := createTempFile("version: \"1.1\"")
		defer os.Remove(filePath)

		config := &Config{}
		err := config.Read(filePath)
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
}
