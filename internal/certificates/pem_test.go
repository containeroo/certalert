package certificates

import (
	"testing"
)

func TestExtractPEMCertificatesInfo(t *testing.T) {
	testCases := []testCase{
		{
			Name:     "Test PEM which no subject",
			FilePath: "../../tests/certs/pem/no_subject.crt",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "1",
					Epoch:   1723889513,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test PEM certificate with password",
			FilePath: "../../tests/certs/pem/with_password.pem",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "with_password",
					Epoch:   1722926986,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test PEM certificate witch is broken",
			FilePath: "../../tests/certs/pem/broken.pem",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode any certificate in 'TestCert'",
		},
		{
			Name:     "Test PEM certificate with wrong password",
			FilePath: "../../tests/certs/pem/with_password.pem",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "with_password",
					Epoch:   1722926986,
					Type:    "pem",
				},
			},
			ExpectedError: "", // no error expected since it it not necessary to decrypt the private key to parse the certificate
		},
		{
			Name:     "Test PEM certificate with chain",
			FilePath: "../../tests/certs/pem/chain.pem",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "final",
					Epoch:   1722926985,
					Type:    "pem",
				},
				{
					Name:    "TestCert",
					Subject: "intermediate",
					Epoch:   1722926986,
					Type:    "pem",
				},
				{
					Name:    "TestCert",
					Subject: "root",
					Epoch:   1722926986,
					Type:    "pem",
				},
			},
			ExpectedError: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if err := unitTest(tc, t, ExtractPEMCertificatesInfo); err != nil {
				t.Error(err)
			}
		})
	}
}
