package certificates

import (
	"testing"
)

func TestExtractJKSCertificatesInfo(t *testing.T) {
	testCases := []testCase{
		{
			Name: "Test JKS certificate - broken",
			Cert: Certificate{
				Name:     "TestCert",
				Path:     "../../tests/certs/jks/broken.jks",
				Password: "password",
			},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to load JKS file 'TestCert': got invalid magic",
		},
		{
			Name: "Test JKS certificate - valid",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
				Path:     "../../tests/certs/jks/regular.jks",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973250,
					Subject: "CN=regular,OU=MyOrganization,O=MyCompany,L=MyCity,ST=MyState,C=MyCountry",
				},
			},
			ExpectedError: "",
		},
		{
			Name: "Test JKS certificate - valid chain",
			Cert: Certificate{
				Name:     "TestCert",
				Password: "password",
				Path:     "../../tests/certs/jks/chain.jks",
			},
			ExpectedResults: []CertificateInfo{
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973251,
					Subject: "CN=root",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973256,
					Subject: "CN=intermediate",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973256,
					Subject: "CN=leaf",
				},
				{
					Name:    "TestCert",
					Type:    "jks",
					Epoch:   1723973251,
					Subject: "CN=chain,OU=MyOrganization,O=MyCompany,L=MyCity,ST=MyState,C=MyCountry",
				},
			},
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if err := runExtractCertificateUnitTest(tc, t, ExtractJKSCertificatesInfo); err != nil {
				t.Error(err)
			}
		})
	}
}
