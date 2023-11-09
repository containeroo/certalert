package certificates

import (
	"testing"
)

func TestExtractJKSCertificatesInfo(t *testing.T) {
	testCases := []testCase{
		{
			Name:     "Test JKS certificate - broken",
			FilePath: "../../tests/certs/jks/broken.jks",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to load JKS file 'TestCert': got invalid magic",
		},
		{
			Name:     "Test JKS certificate - valid",
			FilePath: "../../tests/certs/jks/regular.jks",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973250,
					Subject: "regular",
				},
			},
			ExpectedError: "",
		},
		{
			Name:     "Test JKS certificate - valid chain",
			FilePath: "../../tests/certs/jks/chain.jks",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973251,
					Subject: "root",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973256,
					Subject: "intermediate",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973256,
					Subject: "leaf",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973251,
					Subject: "chain",
				},
			},
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if err := unitTest(tc, t, ExtractJKSCertificatesInfo); err != nil {
				t.Error(err)
			}
		})
	}
}
