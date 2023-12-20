package config

import (
	"certalert/internal/certificates"
	"certalert/internal/utils"
	"fmt"
	"os"
	"sort"
	"testing"
)

// setEnvVars sets all environment variables defined in the given map.
func setEnvVars(envs map[string]string) {
	for key, value := range envs {
		os.Setenv(key, value)
	}
}

// unsetEnvVars unsets all environment variables defined in the given map.
func unsetEnvVars(envs map[string]string) {
	for key := range envs {
		os.Unsetenv(key)
	}
}

// createTempFile creates a temporary file with the given content and returns the file name.
func createTempFile(content string, t *testing.T) string {
	tempFile, err := os.CreateTemp("", "certalert")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer tempFile.Close()

	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	return tempFile.Name()
}

func TestParseCertificatesConfig(t *testing.T) {
	envs := map[string]string{
		"BASIC_PASSWORD": "password",
		"BEARER_TOKEN":   "token",
	}
	passwordFileName := createTempFile("password", t)
	sortedFileExtensions := utils.ExtractMapKeys(certificates.FileExtensionsToType)
	sort.Strings(sortedFileExtensions)

	assertError := func(t *testing.T, expectedError string, actualError error) {
		if expectedError == "" && actualError != nil {
			t.Errorf("Unexpected error: %v", actualError)
		}
		if expectedError != "" {
			if actualError == nil {
				t.Errorf("Expected error '%s', but got none", expectedError)
			} else if actualError.Error() != expectedError {
				t.Errorf("Expected error '%s', but got '%s'", expectedError, actualError.Error())
			}
		}
	}

	t.Run("cert path not defined (FailOnError: false)", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:    "test_cert",
					Enabled: utils.BoolPtr(true),
					Type:    "pem",
				},
			},
			FailOnError: false,
		}
		expectedError := ""

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("cert path not defined", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:    "test_cert",
					Enabled: utils.BoolPtr(true),
					Type:    "pem",
				},
			},
			FailOnError: true,
		}
		expectedError := "Certificate 'test_cert' has no 'path' defined."

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("cert type invalid", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:    "test_cert",
					Enabled: utils.BoolPtr(true),
					Path:    "../../tests/certs/pem/without_password.pem",
					Type:    "invalid",
				},
			},
			FailOnError: true,
		}
		expectedError := fmt.Sprintf("Certificate 'test_cert' has an invalid type 'invalid'. Must be one of %s.", certificates.FileExtensionsTypesSorted)

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("success", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:     "test_cert",
					Path:     "../../tests/certs/pem/chain.pem",
					Type:     "pem",
					Password: fmt.Sprintf("file:%s", passwordFileName),
					Enabled:  utils.BoolPtr(true),
				},
			},
			FailOnError: true,
		}
		expectedError := ""

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("disable cert", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:     "test_cert",
					Path:     "../../tests/certs/pem/chain.pem",
					Enabled:  utils.BoolPtr(false),
					Type:     "pem",
					Password: fmt.Sprintf("file:%s", passwordFileName),
				},
			},
			FailOnError: true,
		}
		expectedError := ""

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("missing_file_var", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:     "test_cert",
					Enabled:  utils.BoolPtr(true),
					Path:     "../../tests/certs/p12/root.p12",
					Type:     "pem",
					Password: "file:INVALID_FILE",
				},
			},
			FailOnError: true,
		}
		expectedError := "Certifacate 'test_cert' has a non resolvable 'password'. Failed to open file 'INVALID_FILE'. open INVALID_FILE: no such file or directory"

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("cert name not defined", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Enabled: utils.BoolPtr(true),
					Path:    "../../tests/certs/pem/final.pem",
				},
			},
			FailOnError: true,
		}
		expectedError := ""

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("cert path not accessible", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:    "test_cert",
					Enabled: utils.BoolPtr(true),
					Path:    "/invalid/path",
				},
			},
			FailOnError: true,
		}
		expectedError := "Certificate 'test_cert' is not accessible. File does not exist: /invalid/path"

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("cert type not defined", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:    "test_cert",
					Enabled: utils.BoolPtr(true),
					Path:    "../../tests/certs/p12/no_extension",
				},
			},
			FailOnError: true,
		}
		expectedError := "Certificate 'test_cert' has no 'type' defined and is missing a file extension."

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("cert type guessed invalid", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:    "test_cert",
					Enabled: utils.BoolPtr(true),
					Path:    "../../tests/certs/pem/cert.invalid",
				},
			},
			FailOnError: true,
		}
		expectedError := "Certificate 'test_cert' has no 'type' defined. Type can't be inferred due to unclear file extension (.invalid)."

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("cert type guessed p12", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:    "test_cert",
					Enabled: utils.BoolPtr(true),
					Path:    "../../tests/certs/p12/chain.p12",
				},
			},
			FailOnError: true,
		}
		expectedError := ""

		setEnvVars(envs)
		err := config.parseCertificatesConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})
}

func TestParsePushgatewayConfig(t *testing.T) {
	envs := map[string]string{
		"PUSHGATEWAY_ADDRESS": "http://localhost:9091",
		"PUSHGATEWAY_JOB":     "certalert",
		"EMPTY_VAR":           "",
	}

	assertError := func(t *testing.T, expectedError string, actualError error) {
		if expectedError == "" && actualError != nil {
			t.Errorf("Unexpected error: %v", actualError)
		}
		if expectedError != "" {
			if actualError == nil {
				t.Errorf("Expected error '%s', but got none", expectedError)
			} else if actualError.Error() != expectedError {
				t.Errorf("Expected error '%s', but got '%s'", expectedError, actualError.Error())
			}
		}
	}

	t.Run("Parse pushgateway missing pushgatway address (FailOnError: false)", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "env:EMPTY_VAR",
			},
			FailOnError: false,
		}
		expectedError := ""

		setEnvVars(envs)
		err := config.parsePushgatewayConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("Pusghateway address VAR is emtpy", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "env:EMPTY_VAR",
			},
			FailOnError: true,
		}
		expectedError := "Pushgateway address was resolved to empty."

		setEnvVars(envs)
		err := config.parsePushgatewayConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("Pushgateway address parsing error", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "http:/localhost:9091",
			},
			FailOnError: true,
		}
		expectedError := "Invalid pushgateway address 'http:/localhost:9091'."

		setEnvVars(envs)
		err := config.parsePushgatewayConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("Pushgateway address no scheme error", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "localhost:9091",
			},
			FailOnError: true,
		}
		expectedError := "Invalid pushgateway address 'localhost:9091'."

		setEnvVars(envs)
		err := config.parsePushgatewayConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("Auth error", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "http://localhost:9091",
				Auth: Auth{
					Basic: &Basic{
						Password: "env:BASIC_PASSWORD",
					},
					Bearer: &Bearer{
						Token: "env:BEARER_TOKEN",
					},
				},
			},
			FailOnError: true,
		}
		expectedError := "Both 'auth.basic' and 'auth.bearer' are defined."

		setEnvVars(envs)
		err := config.parsePushgatewayConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("missing_env_var", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "env:INVALID_ENV",
			},
			FailOnError: true,
		}
		expectedError := "Failed to resolve address for pushgateway. Environment variable 'INVALID_ENV' not found."

		setEnvVars(envs)
		err := config.parsePushgatewayConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("Basic Auth: Fail to resolve password", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "http://localhost:9091",
				Auth: Auth{
					Basic: &Basic{
						Password: "env:INVALID_ENV",
					},
				},
			},
			FailOnError: true,
		}
		expectedError := "Failed to resolve basic auth password for pushgateway. Environment variable 'INVALID_ENV' not found."

		setEnvVars(envs)
		err := config.parsePushgatewayConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("Bearer Auth: Fail to resolve token", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "http://localhost:9091",
				Auth: Auth{
					Bearer: &Bearer{
						Token: "env:INVALID_ENV",
					},
				},
			},
			FailOnError: true,
		}
		expectedError := "Failed to resolve bearer token for pushgateway. Environment variable 'INVALID_ENV' not found."

		setEnvVars(envs)
		err := config.parsePushgatewayConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})

	t.Run("JobName: Fail to resolve job name", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "http://localhost:9091",
				Job:     "env:INVALID_ENV",
			},
			FailOnError: true,
		}
		expectedError := "Failed to resolve jobName for pushgateway. Environment variable 'INVALID_ENV' not found."

		setEnvVars(envs)
		err := config.parsePushgatewayConfig()
		unsetEnvVars(envs)

		assertError(t, expectedError, err)
	})
}

func TestParse(t *testing.T) {
	// Set up environment variables for testing
	envs := map[string]string{
		"PUSHGATEWAY_ADDRESS": "http://localhost:9091",
		"PUSHGATEWAY_JOB":     "certalert",
		"EMPTY_VAR":           "",
	}

	assertError := func(t *testing.T, err error, expectedError interface{}) {
		switch expectedError := expectedError.(type) {
		case nil:
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		case string:
			if err == nil {
				t.Errorf("Expected error '%v', but got nil", expectedError)
			} else if err.Error() != expectedError {
				t.Errorf("Expected error '%v', but got '%v'", expectedError, err)
			}
		default:
			t.Errorf("Unsupported type for expectedError: %T", expectedError)
		}
	}

	t.Run("Pushgateway error", func(t *testing.T) {
		config := &Config{
			Pushgateway: Pushgateway{
				Address: "env:INVALID_ENV",
			},
			FailOnError: true,
		}

		setEnvVars(envs)
		err := config.Parse()
		unsetEnvVars(envs)

		assertError(t, err, "Failed to resolve address for pushgateway. Environment variable 'INVALID_ENV' not found.")
	})

	t.Run("Certificate error", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name: "test_cert",
				},
			},
			Pushgateway: Pushgateway{},
			FailOnError: true,
		}

		setEnvVars(envs)
		err := config.Parse()
		unsetEnvVars(envs)

		assertError(t, err, "Certificate 'test_cert' has no 'path' defined.")
	})

	t.Run("Simple config (success)", func(t *testing.T) {
		config := &Config{
			Certs: []certificates.Certificate{
				{
					Name:    "test_cert",
					Enabled: utils.BoolPtr(true),
					Path:    "../../tests/certs/pem/chain.pem",
					Type:    "pem",
				},
			},
			Pushgateway: Pushgateway{
				Address: "http://localhost:9091",
				Job:     "certalert",
			},
			FailOnError: true,
			Server: Server{
				ListenAddress: "localhost:8080",
			},
		}

		setEnvVars(envs)
		err := config.Parse()
		unsetEnvVars(envs)

		assertError(t, err, nil)
	})
}
