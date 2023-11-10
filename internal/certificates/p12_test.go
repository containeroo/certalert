package certificates

import (
	"testing"
)

func TestExtractP12CertificatesInfo(t *testing.T) {
	testCases := []testCase{
		{
			Name: "Test JKS certificate - JKS with pkcs12",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
				Path:     "../../tests/certs/jks/pkcs12.jks",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Epoch:   1724097113,
					Subject: "CN=pkcs12",
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
		{
			Name: "Test P12 certificate with no password",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "",
				Path:     "../../tests/certs/p12/without_password.p12",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "CN=without_password",
					Epoch:   1722925468,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
		{
			Name: "Test P12 certificate with password",
			Cert: Certificate{
				Name:     "TestCert",
				Path:     "../../tests/certs/p12/with_password.p12",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "CN=with_password",
					Epoch:   1722925469,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
		{
			Name: "Test P12 certificate with wrong password",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "wrong",
				Path:     "../../tests/certs/p12/with_password.p12",
			},
			ExpectedResults: []CertificateInfo{{}},
			ExpectedError:   "Failed to decode P12 file 'TestCert': pkcs12: decryption password incorrect",
		},
		{
			Name: "Test P12 certificate with is broken",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "",
				Path:     "../../tests/certs/p12/broken.p12",
			},
			ExpectedResults: []CertificateInfo{{}},
			ExpectedError:   "Failed to decode P12 file 'TestCert': pkcs12: error reading P12 data: asn1: structure error: tags don't match (16 vs {class:1 tag:2 length:114 isCompound:true}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} pfxPdu @2",
		},
		{
			Name: "Test P12 certificate without subject",
			Cert: Certificate{
				Name:     "TestCert",
				Path:     "../../tests/certs/p12/empty_subject.p12",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "O=Internet Widgits Pty Ltd,ST=Some-State,C=AU",
					Epoch:   1722925468,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
		{
			Name: "Test P12 certificate with chain",
			Cert: Certificate{
				Name:     "TestCert",
				Path:     "../../tests/certs/p12/chain.p12",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "CN=final",
					Epoch:   1722925468,
					Type:    "p12",
				},
				{
					Name:    "TestCert",
					Subject: "CN=intermediate",
					Epoch:   1722925468,
					Type:    "p12",
				},
				{
					Name:    "TestCert",
					Subject: "CN=root",
					Epoch:   1722925468,
					Type:    "p12",
				},
				{
					Name:    "TestCert",
					Subject: "CN=final",
					Epoch:   1722925468,
					Type:    "p12",
				},
			},
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if err := runExtractCertificateUnitTest(tc, t, ExtractP12CertificatesInfo); err != nil {
				t.Error(err)
			}
		})
	}
}
