package config

import "testing"

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
