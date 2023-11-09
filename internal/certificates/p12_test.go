package certificates

import (
	"testing"
)

func TestExtractP12CertificatesInfo(t *testing.T) {
	testCases := []testCase{
		{
			Name:     "Test JKS certificate - JKS with pkcs12",
			FilePath: "../../tests/certs/jks/pkcs12.jks",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "p12",
					Epoch:   1724097113,
					Subject: "pkcs12",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P12 certificate with no password",
			FilePath: "../../tests/certs/p12/without_password.p12",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "without_password",
					Epoch:   1722925468,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P12 certificate with password",
			FilePath: "../../tests/certs/p12/with_password.p12",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "with_password",
					Epoch:   1722925469,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P12 certificate with wrong password",
			FilePath: "../../tests/certs/p12/with_password.p12",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "wrong",
			},
			ExpectedResults: []CertificateInfo{{}},
			ExpectedError:   "Failed to decode P12 file 'TestCert': pkcs12: decryption password incorrect",
		},
		{
			Name:     "Test P12 certificate with is broken",
			FilePath: "../../tests/certs/p12/broken.p12",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "",
			},
			ExpectedResults: []CertificateInfo{{}},
			ExpectedError:   "Failed to decode P12 file 'TestCert': pkcs12: error reading P12 data: asn1: structure error: tags don't match (16 vs {class:1 tag:2 length:114 isCompound:true}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} pfxPdu @2",
		},
		{
			Name:     "Test P12 certificate without subject",
			FilePath: "../../tests/certs/p12/empty_subject.p12",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "1",
					Epoch:   1722925468,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P12 certificate with chain",
			FilePath: "../../tests/certs/p12/chain.p12",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "final",
					Epoch:   1722925468,
					Type:    "p12",
				},
				{
					Name:    "TestCert",
					Subject: "intermediate",
					Epoch:   1722925468,
					Type:    "p12",
				},
				{
					Name:    "TestCert",
					Subject: "root",
					Epoch:   1722925468,
					Type:    "p12",
				},
				{
					Name:    "TestCert",
					Subject: "final",
					Epoch:   1722925468,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if err := unitTest(tc, t, ExtractP12CertificatesInfo); err != nil {
				t.Error(err)
			}
		})
	}
}
