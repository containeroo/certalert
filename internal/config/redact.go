package config

import (
	"certalert/internal/utils"
	"strings"
)

// redactVariable redacts sensitive data from the config if it is not prefixed with env: or file:
func redactVariable(s string) string {
	if strings.HasPrefix(s, "env:") || strings.HasPrefix(s, "file:") {
		return s
	}
	return "<REDACTED>"
}

// RedactConfig redacts sensitive data from the config
func RedactConfig(config *Config) error {

	if utils.HasKey(config.Pushgateway, "Address") {
		config.Pushgateway.Address = redactVariable(config.Pushgateway.Address)
	}

	if utils.HasKey(config.Pushgateway, "Basic.Username") {
		config.Pushgateway.Auth.Basic.Username = redactVariable(config.Pushgateway.Auth.Basic.Username)
	}

	if utils.HasKey(config.Pushgateway, "Basic.Password") {
		config.Pushgateway.Auth.Basic.Password = redactVariable(config.Pushgateway.Auth.Basic.Password)
	}

	if utils.HasKey(config.Pushgateway, "Bearer.Token") {
		config.Pushgateway.Auth.Bearer.Token = redactVariable(config.Pushgateway.Auth.Bearer.Token)
	}

	for idx, cert := range config.Certs {
		config.Certs[idx].Password = redactVariable(cert.Password)
	}

	return nil
}
