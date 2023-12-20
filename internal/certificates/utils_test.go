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

	t.Run("TestCert1", func(t *testing.T) {
		want := &certs[0]
		got, err := GetCertificateByName("TestCert1", certs)
		assert.Equal(t, want, got)
		assert.Nil(t, err)
	})

	t.Run("TestCert3", func(t *testing.T) {
		want := (*Certificate)(nil)
		got, err := GetCertificateByName("TestCert3", certs)
		assert.Equal(t, want, got)
		assert.Equal(t, "Certificate 'TestCert3' not found", err.Error())
	})
}

func TestHandleFailOnError(t *testing.T) {
	var certInfoList []CertificateInfo

	t.Run("failOnError is true", func(t *testing.T) {
		err := handleFailOnError(&certInfoList, "certName", "certType", "An error occurred", true)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		} else if err.Error() != "An error occurred" {
			t.Fatalf("Unexpected error message: %v", err)
		}
	})

	t.Run("failOnError is false", func(t *testing.T) {
		err := handleFailOnError(&certInfoList, "certName", "certType", "Another error occurred", false)
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

	t.Run("Certificate exists in slice", func(t *testing.T) {
		result := certExistsInSlice(cert1, []CertificateInfo{cert1, cert2})
		assert.True(t, result)
	})

	t.Run("Certificate does not exist in slice", func(t *testing.T) {
		result := certExistsInSlice(cert1, []CertificateInfo{cert2, cert3})
		assert.False(t, result)
	})

	t.Run("Certificate with similar but not identical properties", func(t *testing.T) {
		result := certExistsInSlice(CertificateInfo{Name: "cert1", Subject: "subject2", Type: "type2"}, []CertificateInfo{cert1, cert2})
		assert.False(t, result)
	})

	t.Run("Slice is empty", func(t *testing.T) {
		result := certExistsInSlice(cert1, []CertificateInfo{})
		assert.False(t, result)
	})

	t.Run("Slice contains multiple identical certificates", func(t *testing.T) {
		result := certExistsInSlice(cert1, []CertificateInfo{cert1, cert1})
		assert.True(t, result)
	})
}
