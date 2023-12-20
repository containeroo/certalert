package certificates

import (
	"testing"
)

func TestExtractTrustStoreCertificatesInfo(t *testing.T) {
	t.Run("Test TrustStore certificate - broken (FailOnError=true)", func(t *testing.T) {
		tc := testCase{
			Name: "Test TrustStore certificate - broken (FailOnError=true)",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
				Path:     "../../tests/certs/truststore/broken.jks",
			},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode P12 file 'TestCert': pkcs12: error reading P12 data: asn1: structure error: tags don't match (16 vs {class:1 tag:2 length:114 isCompound:true}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} pfxPdu @2",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractTrustStoreCertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test TrustStore certificate - valid", func(t *testing.T) {
		tc := testCase{
			Name: "Test TrustStore certificate - valid",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
				Path:     "../../tests/certs/truststore/regular.jks",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "truststore",
					Epoch:   1722692115,
					Subject: "CN=regular",
				},
			},
			ExpectedError: "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractTrustStoreCertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test TrustStore certificate - valid chain", func(t *testing.T) {
		tc := testCase{
			Name: "Test TrustStore certificate - valid chain",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
				Path:     "../../tests/certs/truststore/chain.jks",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "truststore",
					Epoch:   1722692127,
					Subject: "CN=root",
				},
				{
					Name:    "TestCert",
					Type:    "truststore",
					Epoch:   1722692128,
					Subject: "CN=intermediate",
				},
				{
					Name:    "TestCert",
					Type:    "truststore",
					Epoch:   1722692131,
					Subject: "CN=regular",
				},
			},
			ExpectedError: "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractTrustStoreCertificatesInfo); err != nil {
			t.Error(err)
		}
	})
}
