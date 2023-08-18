package certificates

import (
	"os"
	"testing"
)

func TestExtractP7CertificatesInfo(t *testing.T) {
	type testCase struct {
		Name            string
		FilePath        string
		ExpectedResults []CertificateInfo
		ExpectedError   string
	}

	testCases := []testCase{
		{
			Name:     "Test P7B certificate",
			FilePath: "../../tests/certs/p7/cert1.p7b",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "cert1",
					Epoch:   1723889513,
					Type:    "p7b",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P7B certificate",
			FilePath: "../../tests/certs/p7/cert2.p7b",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "cert2",
					Epoch:   1723889513,
					Type:    "p7b",
				},
			},
			ExpectedError: "",
		},
		{
			Name:            "Test P7B certificate which is broken",
			FilePath:        "../../tests/certs/p7/broken.p7b",
			ExpectedResults: []CertificateInfo{{}},
			ExpectedError:   "Failed to decode any certificate in 'TestCert'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			certData, err := os.ReadFile(tc.FilePath)
			if err != nil {
				t.Fatalf("Failed to read certificate file '%s': %v", tc.Name, err)
			}
			certs, err := ExtractP7CertificatesInfo("TestCert", certData, "", true) // Assuming this function works similarly to the P12 version

			if tc.ExpectedError == "" && err != nil {
				t.Errorf("Test case '%s': unexpected error: %v", tc.Name, err)
			}
			if tc.ExpectedError != "" {
				if err == nil {
					t.Errorf("Test case '%s': expected error, but got none", tc.Name)
				} else if err.Error() != tc.ExpectedError {
					t.Errorf("Test case '%s': expected error '%s', but got '%s'", tc.Name, tc.ExpectedError, err.Error())
				}
				return
			}

			if len(certs) != len(tc.ExpectedResults) {
				t.Errorf("Expected %d certificate, got %d", len(tc.ExpectedResults), len(certs))
			}

			for _, expectedCert := range tc.ExpectedResults {
				var extractedCert CertificateInfo
				for _, cert := range certs {
					if cert.Name == expectedCert.Name &&
						cert.Subject == expectedCert.Subject &&
						cert.Type == expectedCert.Type {
						extractedCert = cert
						break
					}
				}
				if extractedCert == (CertificateInfo{}) {
					t.Errorf("Expected cert %v not found", expectedCert)
				}

				if extractedCert != expectedCert {
					t.Errorf("Expected cert %v, got %v", expectedCert, extractedCert)
				}
			}
		})
	}
}
