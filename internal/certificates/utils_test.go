package certificates

import (
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
