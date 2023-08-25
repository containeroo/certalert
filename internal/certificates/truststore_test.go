package certificates

import (
	"os"
	"testing"
)

func TestExtractTrustStoreCertificatesInfo(t *testing.T) {
	// Define a structure for test case
	type testCase struct {
		Name            string
		FilePath        string
		Password        string
		ExpectedResults []CertificateInfo
		ExpectedError   string
		FailOnError     bool
	}

	// Define test cases
	testCases := []testCase{
		{
			Name:            "Test TrustStore certificate - broken (FailOnError=true)",
			FilePath:        "../../tests/certs/truststore/broken.jks",
			Password:        "password",
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode P12 file 'TestCert': pkcs12: error reading P12 data: asn1: structure error: tags don't match (16 vs {class:1 tag:2 length:114 isCompound:true}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} pfxPdu @2",
			FailOnError:     true,
		},
		{
			Name:     "Test TrustStore certificate - valid",
			FilePath: "../../tests/certs/truststore/regular.jks",
			Password: "password",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "truststore",
					Epoch:   1722692115,
					Subject: "regular",
				},
			},
			ExpectedError: "",
			FailOnError:   true,
		},
		{
			Name:     "Test TrustStore certificate - valid chain",
			FilePath: "../../tests/certs/truststore/chain.jks",
			Password: "password",
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "truststore",
					Epoch:   1722692127,
					Subject: "root",
				},
				{
					Name:    "TestCert",
					Type:    "truststore",
					Epoch:   1722692128,
					Subject: "intermediate",
				},
				{
					Name:    "TestCert",
					Type:    "truststore",
					Epoch:   1722692131,
					Subject: "regular",
				},
			},
			ExpectedError: "",
			FailOnError:   true,
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
			certs, err := ExtractTrustStoreCertificatesInfo("TestCert", certData, tc.Password, tc.FailOnError)

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
