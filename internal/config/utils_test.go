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
