package certificates

import (
	"testing"
)

func TestExtractPEMCertificatesInfo(t *testing.T) {
	testCases := []testCase{
		{
			Name: "Test PEM which no subject",
			Cert: Certificate{
				Name: "TestCert",
				Path: "../../tests/certs/pem/no_subject.crt",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "Certificate 1",
					Epoch:   1723889513,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
		{
			Name: "Test PEM certificate with password",
			Cert: Certificate{
				Name: "TestCert",
				Path: "../../tests/certs/pem/with_password.pem",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "CN=with_password",
					Epoch:   1722926986,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
		{
			Name: "Test PEM certificate witch is broken",
			Cert: Certificate{
				Name: "TestCert",
				Path: "../../tests/certs/pem/broken.pem",
			},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode any certificate in 'TestCert'",
		},
		{
			Name: "Test PEM certificate with wrong password",
			Cert: Certificate{
				Name: "TestCert",
				Path: "../../tests/certs/pem/with_password.pem",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "CN=with_password",
					Epoch:   1722926986,
					Type:    "pem",
				},
			},
			ExpectedError: "", // no error expected since it it not necessary to decrypt the private key to parse the certificate
		},
		{
			Name: "Test PEM certificate with chain",
			Cert: Certificate{
				Name: "TestCert",
				Path: "../../tests/certs/pem/chain.pem",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "CN=final",
					Epoch:   1722926985,
					Type:    "pem",
				},
				{
					Name:    "TestCert",
					Subject: "CN=intermediate",
					Epoch:   1722926986,
					Type:    "pem",
				},
				{
					Name:    "TestCert",
					Subject: "CN=root",
					Epoch:   1722926986,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if err := runExtractCertificateUnitTest(tc, t, ExtractPEMCertificatesInfo); err != nil {
				t.Error(err)
			}
		})
	}
}
