package certificates

import (
	"os"
	"testing"
)

func TestExtractPEMCertificatesInfo(t *testing.T) {
	// Define a structure for test case
	type testCase struct {
		Name            string
		FilePath        string
		ExpectedResults []CertificateInfo
		ExpectedError   string
	}

	// Define test cases
	testCases := []testCase{
		{
			Name:     "Test PEM which no subject",
			FilePath: "../../tests/certs/pem/no_subject.crt",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "1",
					Epoch:   1723889513,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test PEM certificate with password",
			FilePath: "../../tests/certs/pem/with_password.pem",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "with_password",
					Epoch:   1722926986,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
		{
			Name:            "Test PEM certificate witch is broken",
			FilePath:        "../../tests/certs/pem/broken.pem",
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode any certificate in 'TestCert'",
		},
		{
			Name:     "Test PEM certificate with wrong password",
			FilePath: "../../tests/certs/pem/with_password.pem",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "with_password",
					Epoch:   1722926986,
					Type:    "pem",
				},
			},
			ExpectedError: "", // no error expected since it it not necessary to decrypt the private key to parse the certificate
		},
		{
			Name:     "Test PEM certificate with chain",
			FilePath: "../../tests/certs/pem/chain.pem",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "final",
					Epoch:   1722926985,
					Type:    "pem",
				},
				{
					Name:    "TestCert",
					Subject: "intermediate",
					Epoch:   1722926986,
					Type:    "pem",
				},
				{
					Name:    "TestCert",
					Subject: "root",
					Epoch:   1722926986,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Run the function under test
			certData, err := os.ReadFile(tc.FilePath)
			if err != nil {
				t.Errorf("Failed to read certificate file '%s': %v", tc.Name, err)
			}
			certs, err := ExtractPEMCertificatesInfo("TestCert", certData, "", true)

			if tc.ExpectedError == "" && err != nil {
				t.Errorf("Test case '%s': unexpected error: %v", tc.Name, err)
			}
			if tc.ExpectedError != "" {
				if err == nil {
					t.Errorf("Test case '%s': expected error, but got none", tc.Name)
				} else if err.Error() != tc.ExpectedError {
					t.Errorf("Test case '%s': expected error '%s', but got '%s'", tc.Name, tc.ExpectedError, err.Error())
				}
				return // error is expected, so we can skip the rest of the test
			}

			// Check the length of the returned slice
			if len(certs) != len(tc.ExpectedResults) {
				t.Errorf("Expected %d certificates, got %d", len(tc.ExpectedResults), len(certs))
				return
			}

			// Check if each certificate in the expected slice exists in the result slice
			for _, expectedCert := range tc.ExpectedResults {
				if !certExistsInSlice(expectedCert, certs) {
					t.Errorf("Expected cert %v not found", expectedCert)
				}
			}

			// Also check the opposite: each certificate in the result slice should exist in the expected slice
			for _, resultCert := range certs {
				if !certExistsInSlice(resultCert, tc.ExpectedResults) {
					t.Errorf("Unexpected cert found: %v", resultCert)
				}
			}
		})
	}
}
