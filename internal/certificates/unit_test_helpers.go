package certificates

import (
	"fmt"
	"os"
	"testing"
)

type testCase struct {
	Name            string
	FilePath        string
	Cert            Certificate
	ExpectedResults []CertificateInfo
	ExpectedError   string
}

func unitTest(tc testCase, t *testing.T, e extractFunction) error {
	certData, err := os.ReadFile(tc.FilePath)
	if err != nil {
		t.Fatalf("Failed to read certificate file '%s'. %v", tc.Name, err)
	}

	certs, err := e(tc.Cert, certData, true)

	if tc.ExpectedError == "" && err != nil {
		return fmt.Errorf("Test case '%s': unexpected error: %v", tc.Name, err)
	}
	if tc.ExpectedError != "" {
		if err == nil {
			return fmt.Errorf("Test case '%s': expected error, but got none", tc.Name)
		} else if err.Error() != tc.ExpectedError {
			return fmt.Errorf("Test case '%s': expected error '%s', but got '%s'", tc.Name, tc.ExpectedError, err.Error())
		}
		return nil // error is expected, so we can skip the rest of the test
	}

	// Check the length of the returned slice
	if len(certs) != len(tc.ExpectedResults) {
		return fmt.Errorf("Expected %d certificates, got %d", len(tc.ExpectedResults), len(certs))
	}

	// Check if each certificate in the expected slice exists in the result slice
	for _, expectedCert := range tc.ExpectedResults {
		if !certExistsInSlice(expectedCert, certs) {
			return fmt.Errorf("Expected cert %v not found", expectedCert)
		}
	}

	// Also check the opposite: each certificate in the result slice should exist in the expected slice
	for _, resultCert := range certs {
		if !certExistsInSlice(resultCert, tc.ExpectedResults) {
			return fmt.Errorf("Unexpected cert found: %v", resultCert)
		}
	}
	return nil
}
