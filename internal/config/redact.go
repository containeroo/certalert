package config

import (
	"certalert/internal/utils"
	"strings"
)

// RedactConfig redacts sensitive data from a config
// This is a very simple implementation that only redacts the following:
// - Pushgateway.Auth.Basic.Username
// - Pushgateway.Auth.Basic.Password
// - Pushgateway.Auth.Bearer.Token
// - Certs.Password
func RedactConfig(config *Config) error {
	if utils.HasStructField(config, "Pushgateway.Auth.Basic.Username") {
		config.Pushgateway.Auth.Basic.Username = redactVariable(config.Pushgateway.Auth.Basic.Username)
	}

	if utils.HasStructField(config, "Pushgateway.Auth.Basic.Password") {
		config.Pushgateway.Auth.Basic.Password = redactVariable(config.Pushgateway.Auth.Basic.Password)
	}

	if utils.HasStructField(config, "Pushgateway.Auth.Bearer.Token") {
		config.Pushgateway.Auth.Bearer.Token = redactVariable(config.Pushgateway.Auth.Bearer.Token)
	}

	for i, cert := range config.Certs {
		if utils.HasStructField(cert, "Password") {
			config.Certs[i].Password = redactVariable(cert.Password)
		}
	}

	return nil
}

// redactVariable redacts sensitive data from the config if it is not prefixed with env: or file:
func redactVariable(s string) string {
	if strings.HasPrefix(s, "env:") || strings.HasPrefix(s, "file:") || s == "" {
		return s
	}
	return "<REDACTED>"
}
