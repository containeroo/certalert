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
				{Name: "ValidCert1", Path: "tests/certs/jks/regular.jks", Type: "testType", Enabled: utils.BoolPtr(true), Valid: utils.BoolPtr(true)},
				{Name: "ValidCert2", Path: "valid2.jks", Type: "testType", Enabled: utils.BoolPtr(true), Valid: utils.BoolPtr(true)},
			},
			ExpectedInfo: []CertificateInfo{
				{Name: "ValidCert1", Type: "testType", Epoch: 123, Subject: "Test Subject"},
				{Name: "ValidCert2", Type: "testType", Epoch: 123, Subject: "Test Subject"},
			},
		},
		{
			Name:        "skips disabled certificates",
			FailOnError: true,
			Certificates: []Certificate{
				{Name: "DisabledCert", Path: "disabled.jks", Type: "testType", Enabled: utils.BoolPtr(false), Valid: utils.BoolPtr(true)},
			},
			ExpectedInfo: []CertificateInfo{},
		},
		{
			Name:        "skips invalid certificates",
			FailOnError: true,
			Certificates: []Certificate{
				{Name: "InvalidCert", Path: "/certs/jks/broken.jks", Type: "testType", Enabled: utils.BoolPtr(true), Valid: utils.BoolPtr(false)},
			},
			ExpectedInfo: []CertificateInfo{},
		},
		{
			Name:        "fails on extraction error",
			FailOnError: true,
			Certificates: []Certificate{
				{Name: "FailCert", Path: "fail", Type: "testType", Enabled: utils.BoolPtr(true), Valid: utils.BoolPtr(true)},
			},
			ExpectedError: "Error extracting certificate information: Extraction failed",
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
				assert.Equal(t, tc.ExpectedInfo, result)
			}
		})
	}
}
