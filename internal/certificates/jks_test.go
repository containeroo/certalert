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
			FilePath:        "../../tests/certs/jks/broken.jks",
			Password:        "password",
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to load JKS file 'TestCert': got invalid magic",
		},
		{
			Name:     "Test JKS certificate - valid",
			FilePath: "../../tests/certs/jks/regular.jks",
			Password: "password",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973250,
					Subject: "CN=regular,OU=MyOrganization,O=MyCompany,L=MyCity,ST=MyState,C=MyCountry",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test JKS certificate - valid chain",
			FilePath: "../../tests/certs/jks/chain.jks",
			Password: "password",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973251,
					Subject: "CN=root",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973256,
					Subject: "CN=intermediate",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973256,
					Subject: "CN=leaf",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973251,
					Subject: "CN=chain,OU=MyOrganization,O=MyCompany,L=MyCity,ST=MyState,C=MyCountry",
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
				t.Errorf("Failed to read certificate file '%s'. %v", tc.Name, err)
			}
			certs, err := ExtractJKSCertificatesInfo("TestCert", certData, tc.Password, true)

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
