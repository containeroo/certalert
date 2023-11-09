package certificates

import (
	"testing"
)

func TestExtractP7CertificatesInfo(t *testing.T) {
	testCases := []testCase{
		{
			Name:     "Test P7B which no subject",
			FilePath: "../../tests/certs/p7/no_subject.p7b",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "1",
					Epoch:   1723889513,
					Type:    "p7",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P7B certificate which no subject",
			FilePath: "../../tests/certs/p7/no_subject.crt",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "1",
					Epoch:   1723889513,
					Type:    "p7",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P7B certificate with regular certificate",
			FilePath: "../../tests/certs/p7/regular.pem",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "regular",
					Epoch:   1723889513,
					Type:    "p7",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P7B certificate with certificate and regular Private Key",
			FilePath: "../../tests/certs/p7/cert_with_pk.p7",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "regular",
					Epoch:   1723889513,
					Type:    "p7",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P7B certificate with unknown PEM block",
			FilePath: "../../tests/certs/p7/message.p7",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode any certificate in 'TestCert'",
		},
		{
			Name:     "Test P7B certificate",
			FilePath: "../../tests/certs/p7/cert1.p7b",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "cert1",
					Epoch:   1723889513,
					Type:    "p7",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P7B certificate",
			FilePath: "../../tests/certs/p7/cert2.p7b",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Subject: "cert2",
					Epoch:   1723889513,
					Type:    "p7",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test P7B certificate which is broken",
			FilePath: "../../tests/certs/p7/broken.p7b",
			Cert: Certificate{
				Name: "TestCert",
			},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode any certificate in 'TestCert'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if err := unitTest(tc, t, ExtractP7CertificatesInfo); err != nil {
				t.Error(err)
			}
		})
	}
}
