package config

import (
	"certalert/internal/certificates"
	"testing"
)

func TestRedactVariable(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"env:SECRET_ENV", "env:SECRET_ENV"},             // Should not redact env: prefixed strings
		{"file:/path/to/secret", "file:/path/to/secret"}, // Should not redact file: prefixed strings
		{"mysecret", "<REDACTED>"},                       // Should redact non-prefixed strings
		{"", "<REDACTED>"},                               // Should redact empty strings
		{"filemysecret", "<REDACTED>"},                   // Should redact strings that contains the word 'file' but not prefixed with 'file:'
		{"envmysecret", "<REDACTED>"},                    // Should redact strings that contains the word 'env' but not prefixed with 'env:'
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := redactVariable(tt.input)
			if actual != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, actual)
			}
		})
	}
}

func TestRedactConfig(t *testing.T) {
	// Create a mock Config object
	config := &Config{}
	config.Pushgateway.Address = "http://example.com"
	config.Pushgateway.Auth.Basic.Username = "username"
	config.Pushgateway.Auth.Basic.Password = "password"
	config.Pushgateway.Auth.Bearer.Token = "token"
	config.Certs = append(config.Certs, certificates.Certificate{
		Name:     "TestCert",
		Password: "password",
	})

	// Run RedactConfig
	err := RedactConfig(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if the sensitive fields are <REDACTED>
	if config.Pushgateway.Auth.Basic.Username != "<REDACTED>" {
		t.Errorf("Basic Username not <REDACTED>")
	}

	if config.Pushgateway.Auth.Basic.Password != "<REDACTED>" {
		t.Errorf("Basic Password not <REDACTED>")
	}

	if config.Pushgateway.Auth.Bearer.Token != "<REDACTED>" {
		t.Errorf("Bearer Token not <REDACTED>")
	}

	for _, cert := range config.Certs {
		if cert.Password != "<REDACTED>" {
			t.Errorf("Cert Password not <REDACTED>")
		}
	}
}
