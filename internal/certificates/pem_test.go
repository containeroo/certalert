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
		Password        string
		ExpectedResults []CertificateInfo
		ExpectedError   string
	}

	// Define test cases
	testCases := []testCase{
		{
			Name:     "Test PEM certificate with no password",
			FilePath: "../../tests/certs/pem/without_password_certificate.pem",
			Password: "",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "without_password",
					Epoch:   1722689423,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test PEM certificate with password",
			FilePath: "../../tests/certs/pem/with_password_certificate.pem",
			Password: "password",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "with_password",
					Epoch:   1722689423,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
		{
			Name:            "Test PEM certificate witch is broken",
			FilePath:        "../../tests/certs/pem/broken_certificate.pem",
			Password:        "",
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode certificate 'TestCert'",
		},
		{
			Name:     "Test PEM certificate with wrong password",
			FilePath: "../../tests/certs/pem/with_password_certificate.pem",
			Password: "wrong",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "with_password",
					Epoch:   1722689423,
					Type:    "pem",
				},
			},
			ExpectedError: "", // no error expected since it it not necessary to decrypt the private key to parse the certificate
		},
		{
			Name:     "Test PEM certificate with chain",
			FilePath: "../../tests/certs/pem/chain_certificate.pem",
			Password: "password",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "chain",
					Epoch:   1722689422,
					Type:    "pem",
				},
				{
					Name:    "TestCert",
					Subject: "chain",
					Epoch:   1722689422,
					Type:    "pem",
				},
				{
					Name:    "TestCert",
					Subject: "chain",
					Epoch:   1722689422,
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
			certs, err := ExtractPEMCertificatesInfo("TestCert", certData, tc.Password)

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
				t.Errorf("Expected %d certificate, got %d", len(tc.ExpectedResults), len(certs))
			}

			// Check the values in the returned certificate
			for _, expectedCert := range tc.ExpectedResults {
				// Find the certificate in the returned slice
				var extractedCert CertificateInfo
				for _, cert := range certs {
					if cert.Name == expectedCert.Name {
						extractedCert = cert
						break
					}
				}
				if extractedCert == (CertificateInfo{}) {
					t.Errorf("Expected cert %v not found", expectedCert)
				}

				// Check the values in the returned certificate
				if extractedCert != expectedCert {
					t.Errorf("Expected cert %v, got %v", expectedCert, extractedCert)
				}
			}
		})
	}
}
