package config

import (
	"certalert/internal/certificates"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	// Sample configuration for testing.
	original := Config{
		Server: Server{
			Hostname: "localhost",
			Port:     8080,
		},
		Pushgateway: Pushgateway{
			Address: "http://pushgateway.address",
			Job:     "testJob",
			Auth: Auth{
				Basic: Basic{
					Username: "user",
					Password: "pass",
				},
				Bearer: Bearer{
					Token: "token",
				},
			},
		},
		Certs: []certificates.Certificate{
			{
				Name:     "cert1",
				Path:     "/path/to/cert1",
				Password: "cert1password",
				Type:     "jks",
			},
			{
				Name:     "cert2",
				Path:     "/path/to/cert2",
				Password: "cert2password",
				Type:     "pem",
			},
		},
	}

	copy := original.DeepCopy()

	// Check that the copy is not the same as the original.
	if &copy == &original {
		t.Errorf("Copy is the same as the original")
	}

	// modify the copy and check that the original is not modified
	copy.Server.Hostname = "modified"
	if original.Server.Hostname == copy.Server.Hostname {
		t.Errorf("Original is modified")
	}
	copy.Certs[0].Name = "modified"
	if original.Certs[0].Name == copy.Certs[0].Name {
		t.Errorf("Original is modified")
	}

	//modify the original and check that the copy is not modified
	original.Server.Hostname = "original"
	if original.Server.Hostname == copy.Server.Hostname {
		t.Errorf("Copy is modified")
	}
	original.Certs[0].Name = "original"
	if original.Certs[0].Name == copy.Certs[0].Name {
		t.Errorf("Copy is modified")
	}

}

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
