package certificates

import (
	"os"
	"testing"
)

func TestExtractP12CertificatesInfo(t *testing.T) {
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
			Name:     "Test P12 certificate with no password",
			FilePath: "../../tests/certs/p12/without_password.p12",
			Password: "",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "without_password",
					Epoch:   1722925468,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P12 certificate with password",
			FilePath: "../../tests/certs/p12/with_password.p12",
			Password: "password",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "with_password",
					Epoch:   1722925469,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
		{
			Name:            "Test P12 certificate with wrong password",
			FilePath:        "../../tests/certs/p12/with_password.p12",
			Password:        "wrong",
			ExpectedResults: []CertificateInfo{{}},
			ExpectedError:   "Failed to decode P12 file 'TestCert': pkcs12: decryption password incorrect",
		},
		{
			Name:            "Test P12 certificate which is broken",
			FilePath:        "../../tests/certs/p12/broken.p12",
			Password:        "",
			ExpectedResults: []CertificateInfo{{}},
			ExpectedError:   "Failed to decode P12 file 'TestCert': pkcs12: error reading P12 data: asn1: structure error: tags don't match (16 vs {class:1 tag:2 length:114 isCompound:true}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} pfxPdu @2",
		},
		{
			Name:     "Test P12 certificate with chain",
			FilePath: "../../tests/certs/p12/chain.p12",
			Password: "password",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "final",
					Epoch:   1722925468,
					Type:    "p12",
				},
				{
					Name:    "TestCert",
					Subject: "intermediate",
					Epoch:   1722925468,
					Type:    "p12",
				},
				{
					Name:    "TestCert",
					Subject: "root",
					Epoch:   1722925468,
					Type:    "p12",
				},
				{
					Name:    "TestCert",
					Subject: "final",
					Epoch:   1722925468,
					Type:    "p12",
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
			certs, err := ExtractP12CertificatesInfo("TestCert", certData, tc.Password)

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

				// Check the values in the returned certificate
				if extractedCert != expectedCert {
					t.Errorf("Expected cert %v, got %v", expectedCert, extractedCert)
				}
			}
		})
	}
}
