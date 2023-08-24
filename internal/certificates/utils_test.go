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
				{
					Name:     "ValidCert1",
					Path:     "../../tests/certs/jks/regular.jks",
					Password: "password",
					Type:     "jks",
					Enabled:  utils.BoolPtr(true),
				},
				{
					Name:     "ValidCert2",
					Path:     "../../tests/certs/p12/with_password.p12",
					Password: "password",
					Type:     "p12",
					Enabled:  utils.BoolPtr(true),
				},
			},
			ExpectedInfo: []CertificateInfo{
				{
					Name:  "ValidCert1",
					Epoch: 1724096931,
					Type:  "jks", Subject: "regular",
				},
				{
					Name:    "ValidCert2",
					Epoch:   1724097374,
					Type:    "p12",
					Subject: "with_password",
				},
			},
		},
		{
			Name:        "skips disabled certificate",
			FailOnError: true,
			Certificates: []Certificate{
				{
					Name:    "DisabledCert",
					Path:    "disabled.jks",
					Type:    "jks",
					Enabled: utils.BoolPtr(false),
				},
			},
			ExpectedInfo: []CertificateInfo(nil),
		},
		{
			Name:        "skips invalid certificate",
			FailOnError: true,
			Certificates: []Certificate{
				{
					Name:    "InvalidCert",
					Path:    "../tests/certs/jks/broken.jks",
					Type:    "jks",
					Enabled: utils.BoolPtr(true),
				},
			},
			ExpectedInfo:  []CertificateInfo(nil),
			ExpectedError: "Failed to read certificate file '../tests/certs/jks/broken.jks': open ../tests/certs/jks/broken.jks: no such file or directory",
		},
		{
			Name:        "fails on extraction error",
			FailOnError: true,
			Certificates: []Certificate{
				{
					Name:    "FailCert",
					Path:    "fail",
					Type:    "jks",
					Enabled: utils.BoolPtr(true),
				},
			},
			ExpectedError: "Failed to read certificate file 'fail': open fail: no such file or directory",
		},
		{
			Name:        "fails on invalid type (failsOnError = true)",
			FailOnError: true,
			Certificates: []Certificate{
				{
					Name:     "InvalidTypeCert",
					Path:     "../../tests/certs/jks/regular.jks",
					Password: "password",
					Type:     "invalid",
					Enabled:  utils.BoolPtr(true),
				},
			},
			ExpectedError: "Unknown certificate type 'invalid'",
		},
		{
			Name:        "fails on invalid password (failsOnError = true)",
			FailOnError: true,
			Certificates: []Certificate{
				{
					Name:     "InvalidPasswordCert",
					Path:     "../../tests/certs/p12/with_password.p12",
					Password: "invalid",
					Type:     "p12",
					Enabled:  utils.BoolPtr(true),
				},
			},
			ExpectedError: "Error extracting certificate information: Failed to decode P12 file 'InvalidPasswordCert': pkcs12: decryption password incorrect",
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

				// Check the length of the returned slice
				if len(result) != len(tc.ExpectedInfo) {
					t.Errorf("Expected %d certificates, got %d", len(tc.ExpectedInfo), len(result))
					return
				}

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

	t.Run("fails on invalid type (failsOnError = false)", func(t *testing.T) {
		certs := []Certificate{
			{
				Name:     "InvalidTypeCert",
				Path:     "../../tests/certs/jks/regular.jks",
				Password: "password",
				Type:     "invalid",
				Enabled:  utils.BoolPtr(true),
			},
		}

		result, err := Process(certs, false)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "InvalidTypeCert", result[0].Name)
		assert.Equal(t, "invalid", result[0].Type)
		assert.Equal(t, "Unknown certificate type 'invalid'", result[0].Error)
	})
}

func TestHandleError(t *testing.T) {
	var certInfoList []CertificateInfo

	t.Run("failOnError is true", func(t *testing.T) {
		err := handleError(&certInfoList, "certName", "certType", "An error occurred", true)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		} else if err.Error() != "An error occurred" {
			t.Fatalf("Unexpected error message: %v", err)
		}
	})

	t.Run("failOnError is false", func(t *testing.T) {
		err := handleError(&certInfoList, "certName", "certType", "Another error occurred", false)
		if err != nil {
			t.Fatalf("Did not expect an error, got: %v", err)
		}

		if len(certInfoList) != 1 {
			t.Fatalf("Expected certInfoList to have 1 entry, got: %d", len(certInfoList))
		}

		certInfo := certInfoList[0]
		if certInfo.Name != "certName" || certInfo.Type != "certType" || certInfo.Error != "Another error occurred" {
			t.Fatalf("Unexpected entry in certInfoList: %+v", certInfo)
		}
	})
}

func TestCertExistsInSlice(t *testing.T) {
	cert1 := CertificateInfo{Name: "cert1", Subject: "subject1", Type: "type1"}
	cert2 := CertificateInfo{Name: "cert2", Subject: "subject2", Type: "type2"}
	cert3 := CertificateInfo{Name: "cert3", Subject: "subject3", Type: "type3"}

	tests := []struct {
		name      string
		cert      CertificateInfo
		certSlice []CertificateInfo
		expected  bool
	}{
		{
			name:      "Certificate exists in slice",
			cert:      cert1,
			certSlice: []CertificateInfo{cert1, cert2},
			expected:  true,
		},
		{
			name:      "Certificate does not exist in slice",
			cert:      cert1,
			certSlice: []CertificateInfo{cert2, cert3},
			expected:  false,
		},
		{
			name:      "Certificate with similar but not identical properties",
			cert:      CertificateInfo{Name: "cert1", Subject: "subject2", Type: "type2"},
			certSlice: []CertificateInfo{cert1, cert2},
			expected:  false,
		},
		{
			name:      "Slice is empty",
			cert:      cert1,
			certSlice: []CertificateInfo{},
			expected:  false,
		},
		{
			name:      "Slice contains multiple identical certificates",
			cert:      cert1,
			certSlice: []CertificateInfo{cert1, cert1},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := certExistsInSlice(tt.cert, tt.certSlice)
			assert.Equal(t, tt.expected, result)
		})
	}
}
