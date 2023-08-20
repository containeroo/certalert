package certificates

import (
	"certalert/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test cases for GetCertificateByName
func TestGetCertificateByName(t *testing.T) {
	certs := []Certificate{
		{Name: "TestCert1"},
		{Name: "TestCert2"},
	}

	tt := []struct {
		name string
		want *Certificate
		err  string
	}{
		{"TestCert1", &certs[0], ""},
		{"TestCert3", nil, "Certificate 'TestCert3' not found"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetCertificateByName(tc.name, certs)
			assert.Equal(t, tc.want, got)
			if err != nil {
				assert.Equal(t, err.Error(), tc.err)
			}
		})
	}
}

func TestProcess(t *testing.T) {
	cases := []struct {
		Name          string
		Certificates  []Certificate
		FailOnError   bool
		ExpectedInfo  []CertificateInfo
		ExpectedError string
	}{
		{
			Name:        "handles valid certificates",
			FailOnError: true,
			Certificates: []Certificate{
				{Name: "ValidCert1", Path: "../../tests/certs/jks/regular.jks", Password: "password", Type: "jks", Enabled: utils.BoolPtr(true), Valid: utils.BoolPtr(true)},
				{Name: "ValidCert2", Path: "../../tests/certs/p12/with_password.p12", Password: "password", Type: "p12", Enabled: utils.BoolPtr(true), Valid: utils.BoolPtr(true)},
			},
			ExpectedInfo: []CertificateInfo{
				{Name: "ValidCert1", Epoch: 1724096931, Type: "jks", Subject: "regular"},
				{Name: "ValidCert2", Epoch: 1724097374, Type: "p12", Subject: "with_password"},
			},
		},
		{
			Name:        "skips disabled certificate",
			FailOnError: true,
			Certificates: []Certificate{
				{Name: "DisabledCert", Path: "disabled.jks", Type: "jks", Enabled: utils.BoolPtr(false), Valid: utils.BoolPtr(true)},
			},
			ExpectedInfo: []CertificateInfo(nil),
		},
		{
			Name:        "skips invalid certificate",
			FailOnError: true,
			Certificates: []Certificate{
				{Name: "InvalidCert", Path: "../tests/certs/jks/broken.jks", Type: "jks", Enabled: utils.BoolPtr(true), Valid: utils.BoolPtr(false)},
			},
			ExpectedInfo: []CertificateInfo(nil),
		},
		{
			Name:        "fails on extraction error",
			FailOnError: true,
			Certificates: []Certificate{
				{Name: "FailCert", Path: "fail", Type: "jks", Enabled: utils.BoolPtr(true), Valid: utils.BoolPtr(true)},
			},
			ExpectedError: "Failed to read certificate file 'fail': open fail: no such file or directory",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := Process(tc.Certificates, tc.FailOnError)
			if tc.ExpectedError != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tc.ExpectedError, err.Error())
			} else {
				assert.Nil(t, err)

				// Check if each certificate in the expected slice exists in the result slice
				for _, expectedCert := range tc.ExpectedInfo {
					if !certExistsInSlice(expectedCert, result) {
						t.Errorf("Expected cert %v not found", expectedCert)
					}
				}

				// Also check the opposite: each certificate in the result slice should exist in the expected slice
				for _, resultCert := range result {
					if !certExistsInSlice(resultCert, tc.ExpectedInfo) {
						t.Errorf("Unexpected cert found: %v", resultCert)
					}
				}

			}
		})
	}
}
