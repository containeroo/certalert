package certificates

import (
	"os"
	"testing"
)

func TestExtractJKSCertificatesInfo(t *testing.T) {
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
			Name:            "Test JKS certificate - broken",
			FilePath:        "../../tests/certs/jks/broken_certificate.jks",
			Password:        "password",
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to load JKS file 'TestCert': got invalid magic",
		},
		{
			Name:     "Test JKS certificate - valid",
			FilePath: "../../tests/certs/jks/regular.jks",
			Password: "changeit",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1722692115,
					Subject: "regular",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test JKS certificate - valid chain",
			FilePath: "../../tests/certs/jks/chain.jks",
			Password: "changeit",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1722692127,
					Subject: "root",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1722692128,
					Subject: "intermediate",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1722692131,
					Subject: "leaf",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1722692126,
					Subject: "chain",
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
			certs, err := ExtractJKSCertificatesInfo("TestCert", certData, tc.Password)

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
				var found bool
				for _, cert := range certs {
					if cert.Name == expectedCert.Name && cert.Subject == expectedCert.Subject && cert.Type == expectedCert.Type {
						extractedCert = cert
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected cert %v not found", expectedCert)
				}
				if expectedCert.Epoch != extractedCert.Epoch {
					t.Errorf("Expected cert %v, got %v", expectedCert, extractedCert)
				}
			}
		})
	}
}
