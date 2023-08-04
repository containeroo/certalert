package config

import (
	"certalert/internal/certificates"
	"fmt"
	"os"
	"strings"
	"testing"
)

func setEnvVars(envs map[string]string) {
	for key, value := range envs {
		os.Setenv(key, value)
	}
}

func unsetEnvVars(envs map[string]string) {
	for key := range envs {
		os.Unsetenv(key)
	}
}

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
		"PUSHGATEWAY_ADDRESS": "http://localhost:9091",
		"PUSHGATEWAY_JOB":     "certalert",
		"BASIC_PASSWORD":      "password",
		"BEARER_TOKEN":        "token",
	}
	passwordFileName := createTempFile("password", t)
	var trueVar = true

	testCases := []struct {
		name          string
		config        *Config
		expectedError string
	}{
		{
			name: "success",
			config: &Config{
				Pushgateway: Pushgateway{
					Address: "env:PUSHGATEWAY_ADDRESS",
					Job:     "env:PUSHGATEWAY_JOB",
				},
				Certs: []certificates.Certificate{
					{
						Name:     "test_cert",
						Path:     "path/to/cert.pem",
						Type:     "pem",
						Password: fmt.Sprintf("file:%s", passwordFileName),
						Enabled:  &trueVar,
					},
				},
			},
			expectedError: "",
		},
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
			},
			expectedError: "Both 'auth.basic' and 'auth.bearer' are defined",
		},
		{
			name: "missing_env_var",
			config: &Config{
				Pushgateway: Pushgateway{
					Address: "env:INVALID_ENV",
				},
			},
			expectedError: "Failed to resolve address for pushgateway: Environment variable 'INVALID_ENV' not found",
		},
		{
			name: "missing_file_var",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:     "test_cert",
						Enabled:  &trueVar,
						Path:     "path/to/cert.pem",
						Type:     "pem",
						Password: "file:INVALID_FILE",
					},
				},
			},
			expectedError: "Certifacate 'test_cert' cannot resolve 'password'. Failed to open file 'INVALID_FILE': open INVALID_FILE: no such file or directory",
		},
		{
			name: "cert name not defined",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Enabled: &trueVar,
						Path:    "path/to/cert-with_special chars.pem",
						Name:    "cert-with-special-chars-pem",
					},
				},
			},
			expectedError: "",
		},
		{
			name: "cert path not defined",
			config: &Config{
				Certs: []certificates.Certificate{
					{

						Name:    "test_cert",
						Enabled: &trueVar,
						Type:    "pem",
					},
				},
			},
			expectedError: "Certificate 'test_cert' has no 'path' defined",
		},
		{
			name: "cert type not defined",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: &trueVar,
						Path:    "path/to/cert",
					},
				},
			},
			expectedError: "Certificate 'test_cert' has no 'type' defined. Type can't be inferred due to the missing file extension.",
		},
		{
			name: "cert type invalid",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: &trueVar,
						Path:    "path/to/cert.invalid",
						Type:    "invalid",
					},
				},
			},
			expectedError: fmt.Sprintf("Certificate 'test_cert' has an invalid 'type'. Must be one of: %s", strings.Join(certificates.ValidTypes, ", ")),
		},
		{
			name: "cert type guessed invalid",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: &trueVar,
						Path:    "path/to/cert.invalid",
					},
				},
			},
			expectedError: "Certificate 'test_cert' has no 'type' defined. Type can't be inferred due to the unclear file extension (.invalid).",
		},
		{
			name: "cert type guessed p12",
			config: &Config{
				Certs: []certificates.Certificate{
					{
						Name:    "test_cert",
						Enabled: &trueVar,
						Path:    "path/to/cert.p12",
					},
				},
			},
			expectedError: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			setEnvVars(envs)

			fmt.Println("test name:", testCase.name)
			err := ParseConfig(testCase.config)

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
