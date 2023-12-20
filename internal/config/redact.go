package config

import (
	"certalert/internal/utils"
	"strings"
)

// RedactConfig redacts sensitive data from a configuration object.
//
// This function redacts sensitive information in the provided Config object, such as:
// - Pushgateway.Auth.Basic.Username
// - Pushgateway.Auth.Basic.Password
// - Pushgateway.Auth.Bearer.Token
// - Certs.Password
//
// Parameters:
//   - config: *Config
//     A pointer to the Config object to be redacted.
//
// Returns:
//   - error
//     An error if redacting the sensitive information fails.
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

// redactVariable redacts sensitive data from a string if it is not prefixed with "env:" or "file:".
//
// This function is used to redact sensitive information from a string, such as passwords, unless the string
// is explicitly marked with the "env:" or "file:" prefix. If the input string is empty or already prefixed,
// it remains unchanged; otherwise, it is redacted.
//
// Parameters:
//   - s: string
//     The string to be redacted.
//
// Returns:
//   - string
//     The redacted or unchanged string.
func redactVariable(s string) string {
	if strings.HasPrefix(s, "env:") || strings.HasPrefix(s, "file:") || s == "" {
		return s
	}
	return "<REDACTED>"
}
