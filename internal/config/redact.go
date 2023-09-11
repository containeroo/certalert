package config

import (
	"certalert/internal/utils"
	"fmt"
	"strings"
)

// RedactConfig redacts sensitive data from a config
// This is a very simple implementation that only redacts the following:
// - Pushgateway.Auth.Basic.Username
// - Pushgateway.Auth.Basic.Password
// - Pushgateway.Auth.Bearer.Token
// - Certs.Password
func RedactConfig(config *Config) error {
	toRedact := []string{
		"Pushgateway.Auth.Basic.Username",
		"Pushgateway.Auth.Basic.Password",
		"Pushgateway.Auth.Bearer.Token",
		"Certs[].Password",
	}

	for _, path := range toRedact {
		if err := utils.UpdateFieldByPath(config, path, redactVariable); err != nil {
			return fmt.Errorf("Failed to redact config: %s", err)
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
