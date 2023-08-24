package config

import (
	"certalert/internal/certificates"
	"certalert/internal/utils"
	"fmt"
	"os"
	"sort"
	"strings"
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

func TestParseConfig(t *testing.T) {
	envs := map[string]string{
		"BASIC_PASSWORD": "password",
		"BEARER_TOKEN":   "token",
	}
	passwordFileName := createTempFile("password", t)
	sortedFileExtensions := utils.ExtractMapKeys(certificates.FileExtensionsToType, true)
	sort.Strings(sortedFileExtensions)

	testCases := []struct {
		name          string
		config        *Config
		expectedError string
	}{
		{
			name: "success",
			config: &Config{
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
			},
			expectedError: "",
		},
		{
			name: "disable cert",
			config: &Config{
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
			},
			expectedError: "",
		},
		{
			name: "missing_file_var",
			config: &Config{
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
			},
			expectedError: "Certifacate 'test_cert' has a non resolvable 'password'. Failed to open file 'INVALID_FILE': open INVALID_FILE: no such file or directory",
		},
		{
			name: "cert name not defined",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Enabled: utils.BoolPtr(true),
						Path:    "../../tests/certs/pem/final.pem",
					},
				},
				FailOnError: true,
			},
			expectedError: "",
		},
		{
			name: "cert path not defined",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: utils.BoolPtr(true),
						Type:    "pem",
					},
				},
				FailOnError: true,
			},
			expectedError: "Certificate 'test_cert' has no 'path' defined",
		},
		{
			name: "cert path not accessible",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: utils.BoolPtr(true),
						Path:    "/invalid/path",
					},
				},
				FailOnError: true,
			},
			expectedError: "Certificate 'test_cert' is not accessible: File does not exist: /invalid/path",
		},
		{
			name: "cert type not defined",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: utils.BoolPtr(true),
						Path:    "../../tests/certs/p12/no_extension",
					},
				},
				FailOnError: true,
			},
			expectedError: "Certificate 'test_cert' has no 'type' defined. Type can't be inferred due to the missing file extension.",
		},
		{
			name: "cert type invalid",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: utils.BoolPtr(true),
						Path:    "../../tests/certs/pem/without_password.pem",
						Type:    "invalid",
					},
				},
				FailOnError: true,
			},
			expectedError: fmt.Sprintf("Certificate 'test_cert' has an invalid 'type'. Must be one of: %s", strings.Join(sortedFileExtensions, ", ")),
		},
		{
			name: "cert type guessed invalid",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: utils.BoolPtr(true),
						Path:    "../../tests/certs/pem/cert.invalid",
					},
				},
				FailOnError: true,
			},
			expectedError: "Certificate 'test_cert' has no 'type' defined. Type can't be inferred due to the unclear file extension (.invalid).",
		},
		{
			name: "cert type guessed p12",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: utils.BoolPtr(true),
						Path:    "../../tests/certs/p12/chain.p12",
					},
				},
				FailOnError: true,
			},
			expectedError: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			setEnvVars(envs)

			err := testCase.config.parseCertificatesConfig()

			unsetEnvVars(envs)

			if testCase.expectedError == "" && err != nil {
				t.Errorf("Test case '%s': unexpected error: %v", testCase.name, err)
			}
			if testCase.expectedError != "" {
				if err == nil {
					t.Errorf("Test case '%s': expected error, but got none", testCase.name)
				} else if err.Error() != testCase.expectedError {
					t.Errorf("Test case '%s': expected error '%s', but got '%s'", testCase.name, testCase.expectedError, err.Error())
				}
			}
			// if reached here, we have no error, so we can continue with the next test case
		})
	}
}

func TestParsePushgatewayConfig(t *testing.T) {
	envs := map[string]string{
		"PUSHGATEWAY_ADDRESS": "http://localhost:9091",
		"PUSHGATEWAY_JOB":     "certalert",
	}

	testCases := []struct {
		name          string
		config        *Config
		expectedError string
	}{
		{
			name: "Auth error",
			config: &Config{
				Pushgateway: Pushgateway{
					Auth: Auth{
						Basic: Basic{
							Password: "env:BASIC_PASSWORD",
						},
						Bearer: Bearer{
							Token: "env:BEARER_TOKEN",
						},
					},
				},
				FailOnError: true,
			},
			expectedError: "Both 'auth.basic' and 'auth.bearer' are defined",
		},
		{
			name: "missing_address",
			config: &Config{
				Pushgateway: Pushgateway{},
				FailOnError: true,
			},
		},
		{
			name: "invalid_address",
			config: &Config{
				Pushgateway: Pushgateway{
					Address: "invalid",
				},
				FailOnError: true,
			},
		},
		{
			name: "missing_env_var",
			config: &Config{
				Pushgateway: Pushgateway{
					Address: "env:INVALID_ENV",
				},
				FailOnError: true,
			},
			expectedError: "Failed to resolve address for pushgateway: Environment variable 'INVALID_ENV' not found.",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			setEnvVars(envs)

			err := testCase.config.parsePushgatewayConfig()

			unsetEnvVars(envs)

			if testCase.expectedError == "" && err != nil {
				t.Errorf("Test case '%s': unexpected error: %v", testCase.name, err)
			}
			if testCase.expectedError != "" {
				if err == nil {
					t.Errorf("Test case '%s': expected error, but got none", testCase.name)
				} else if err.Error() != testCase.expectedError {
					t.Errorf("Test case '%s': expected error '%s', but got '%s'", testCase.name, testCase.expectedError, err.Error())
				}
			}
			// if reached here, we have no error, so we can continue with the next test case
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name          string
		config        *Config
		expectedError error
	}{
		{
			name: "Simple config (success)",
			config: &Config{
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
			},
			expectedError: nil,
		},
		{
			name: "Pushgateway error",
			config: &Config{
				Pushgateway: Pushgateway{
					Address: "env:INVALID_ENV",
				},
				FailOnError: true,
			},
			expectedError: fmt.Errorf("Failed to resolve address for pushgateway: Environment variable 'INVALID_ENV' not found."),
		},
		{
			name: "Certificate error",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name: "test_cert",
					},
				},
				Pushgateway: Pushgateway{},
				FailOnError: true,
			},
			expectedError: fmt.Errorf("Certificate 'test_cert' has no 'path' defined"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Parse()
			if err != nil {
				if tc.expectedError == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("Expected error: '%v', got: '%v'", tc.expectedError, err)
				}
			} else if tc.expectedError != nil {
				t.Errorf("Expected error: %v, got nil", tc.expectedError)
			}
		})
	}
}
