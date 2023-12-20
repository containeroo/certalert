package config

import (
	"certalert/internal/certificates"
	"testing"
)

func TestRedactVariable(t *testing.T) {
	t.Run("does not redact env: prefixed strings", func(t *testing.T) {
		input := "env:SECRET_ENV"
		expected := "env:SECRET_ENV"
		actual := redactVariable(input)
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("does not redact file: prefixed strings", func(t *testing.T) {
		input := "file:/path/to/secret"
		expected := "file:/path/to/secret"
		actual := redactVariable(input)
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("redacts non-prefixed strings", func(t *testing.T) {
		input := "mysecret"
		expected := "<REDACTED>"
		actual := redactVariable(input)
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("does not redact empty strings", func(t *testing.T) {
		input := ""
		expected := ""
		actual := redactVariable(input)
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("redacts strings that contain 'file' but are not prefixed with 'file:'", func(t *testing.T) {
		input := "filemysecret"
		expected := "<REDACTED>"
		actual := redactVariable(input)
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("redacts strings that contain 'env' but are not prefixed with 'env:'", func(t *testing.T) {
		input := "envmysecret"
		expected := "<REDACTED>"
		actual := redactVariable(input)
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})
}

func TestRedactConfig(t *testing.T) {
	// Create a mock Config object
	config := &Config{}
	config.Pushgateway.Address = "http://example.com"

	a := &Auth{
		Basic: &Basic{
			Username: "username",
			Password: "password",
		},
		Bearer: &Bearer{
			Token: "token",
		},
	}

	config.Pushgateway.Auth = *a
	config.Certs = append(config.Certs, certificates.Certificate{
		Name:     "TestCert",
		Password: "password",
	})

	err := RedactConfig(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Run("does not redact non-sensitive fields", func(t *testing.T) {
		if config.Pushgateway.Address != "http://example.com" {
			t.Errorf("Address not http://example.com")
		}
	})

	t.Run("redacts sensitive fields", func(t *testing.T) {
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
	})
}
