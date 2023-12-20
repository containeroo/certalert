package certificates

import (
	"testing"
)

func TestExtractPEMCertificatesInfo(t *testing.T) {
	t.Run("Test PEM which no subject", func(t *testing.T) {
		tc := testCase{
			Name:            "Test PEM which no subject",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/pem/no_subject.crt"},
			ExpectedResults: []CertificateInfo{{Name: "TestCert", Subject: "Certificate 1", Epoch: 1723889513, Type: "pem"}},
			ExpectedError:   "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractPEMCertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test PEM certificate with password", func(t *testing.T) {
		tc := testCase{
			Name:            "Test PEM certificate with password",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/pem/with_password.pem"},
			ExpectedResults: []CertificateInfo{{Name: "TestCert", Subject: "CN=with_password", Epoch: 1722926986, Type: "pem"}},
			ExpectedError:   "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractPEMCertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test PEM certificate witch is broken", func(t *testing.T) {
		tc := testCase{
			Name:            "Test PEM certificate witch is broken",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/pem/broken.pem"},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode any certificate in 'TestCert'",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractPEMCertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test PEM certificate with wrong password", func(t *testing.T) {
		tc := testCase{
			Name:            "Test PEM certificate with wrong password",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/pem/with_password.pem"},
			ExpectedResults: []CertificateInfo{{Name: "TestCert", Subject: "CN=with_password", Epoch: 1722926986, Type: "pem"}},
			ExpectedError:   "", // no error expected since it is not necessary to decrypt the private key to parse the certificate
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractPEMCertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test PEM certificate with chain", func(t *testing.T) {
		tc := testCase{
			Name: "Test PEM certificate with chain",
			Cert: Certificate{Name: "TestCert", Path: "../../tests/certs/pem/chain.pem"},
			ExpectedResults: []CertificateInfo{
				{Name: "TestCert", Subject: "CN=final", Epoch: 1722926985, Type: "pem"},
				{Name: "TestCert", Subject: "CN=intermediate", Epoch: 1722926986, Type: "pem"},
				{Name: "TestCert", Subject: "CN=root", Epoch: 1722926986, Type: "pem"},
			},
			ExpectedError: "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractPEMCertificatesInfo); err != nil {
			t.Error(err)
		}
	})
}
