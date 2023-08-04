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
