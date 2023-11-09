package certificates

import (
	"testing"
)

func TestExtractTrustStoreCertificatesInfo(t *testing.T) {
	testCases := []testCase{
		{
			Name:     "Test TrustStore certificate - broken (FailOnError=true)",
			FilePath: "../../tests/certs/truststore/broken.jks",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode P12 file 'TestCert': pkcs12: error reading P12 data: asn1: structure error: tags don't match (16 vs {class:1 tag:2 length:114 isCompound:true}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} pfxPdu @2",
		},
		{
			Name:     "Test TrustStore certificate - valid",
			FilePath: "../../tests/certs/truststore/regular.jks",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "truststore",
					Epoch:   1722692115,
					Subject: "regular",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test TrustStore certificate - valid chain",
			FilePath: "../../tests/certs/truststore/chain.jks",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
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
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if err := unitTest(tc, t, ExtractTrustStoreCertificatesInfo); err != nil {
				t.Error(err)
			}
		})
	}
}
