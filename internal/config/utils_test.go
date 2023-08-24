package config

import (
	"testing"
)

func TestExtractHostAndPort(t *testing.T) {
	tests := []struct {
		input        string
		expectedHost string
		expectedPort int
		expectedErr  bool
	}{
		{"example.com:8080", "example.com", 8080, false},
		{":1234", "", 1234, false},
		{"localhost:", "", 0, true},
		{"localhost:8080", "localhost", 8080, false},
		{"127.0.0.1:", "", 0, true},
		{"127.0.0.1:8080", "127.0.0.1", 8080, false},
		{"invalid", "", 0, true},
		{"invalid:", "", 0, true},
	}

	for _, test := range tests {
		host, port, err := ExtractHostAndPort(test.input)

		if (err != nil) != test.expectedErr {
			t.Errorf("For %s, expected error: %v, but got: %v", test.input, test.expectedErr, err != nil)
			continue
		}

		if host != test.expectedHost {
			t.Errorf("For %s, expected host: %s, but got: %s", test.input, test.expectedHost, host)
		}

		if port != test.expectedPort {
			t.Errorf("For %s, expected port: %d, but got: %d", test.input, test.expectedPort, port)
		}
	}
}
