package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Read reads the config file
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
