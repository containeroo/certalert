package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ReadConfigFile reads the config file and unmarshals it into the Config struct
func ReadConfigFile(configPath string, destination interface{}) error {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(destination); err != nil {
		return fmt.Errorf("Failed to unmarshal config file: %v", err)
	}

	return nil
}
