package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Read reads the configuration settings from the specified config file.
//
// This method sets the Viper config file path and reads the config file using Viper.
// It then unmarshals the configuration into the provided Config struct using mapstructure,
// with the option to zero out any existing fields.
//
// Parameters:
//   - configPath: string
//     The file path to the configuration file.
//
// Returns:
//   - error
//     An error if reading or unmarshaling the configuration fails.
func (c *Config) Read(configPath string) error {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(c, func(d *mapstructure.DecoderConfig) { d.ZeroFields = true }); err != nil {
		return fmt.Errorf("Failed to unmarshal config file: %v", err)
	}

	return nil
}
