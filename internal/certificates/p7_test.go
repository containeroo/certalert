package certificates

import "testing"

func TestExtractP7CertificatesInfo(t *testing.T) {
	t.Run("Test P7B which no subject", func(t *testing.T) {
		tc := testCase{
			Name:            "Test P7B which no subject",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/p7/no_subject.p7b"},
			ExpectedResults: []CertificateInfo{{Name: "TestCert", Subject: "Certificate 1", Epoch: 1723889513, Type: "p7"}},
			ExpectedError:   "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractP7CertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test P7B certificate which no subject", func(t *testing.T) {
		tc := testCase{
			Name:            "Test P7B certificate which no subject",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/p7/no_subject.crt"},
			ExpectedResults: []CertificateInfo{{Name: "TestCert", Subject: "Certificate 1", Epoch: 1723889513, Type: "p7"}},
			ExpectedError:   "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractP7CertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test P7B certificate with regular certificate", func(t *testing.T) {
		tc := testCase{
			Name:            "Test P7B certificate with regular certificate",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/p7/regular.pem"},
			ExpectedResults: []CertificateInfo{{Name: "TestCert", Subject: "CN=regular", Epoch: 1723889513, Type: "p7"}},
			ExpectedError:   "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractP7CertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test P7B certificate with certificate and regular Private Key", func(t *testing.T) {
		tc := testCase{
			Name:            "Test P7B certificate with certificate and regular Private Key",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/p7/cert_with_pk.p7"},
			ExpectedResults: []CertificateInfo{{Name: "TestCert", Subject: "CN=regular", Epoch: 1723889513, Type: "p7"}},
			ExpectedError:   "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractP7CertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test P7B certificate with unknown PEM block", func(t *testing.T) {
		tc := testCase{
			Name:            "Test P7B certificate with unknown PEM block",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/p7/message.p7"},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode any certificate in 'TestCert'",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractP7CertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test P7B certificate", func(t *testing.T) {
		tc := testCase{
			Name:            "Test P7B certificate",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/p7/cert1.p7b"},
			ExpectedResults: []CertificateInfo{{Name: "TestCert", Subject: "CN=cert1", Epoch: 1723889513, Type: "p7"}},
			ExpectedError:   "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractP7CertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test P7B certificate", func(t *testing.T) {
		tc := testCase{
			Name:            "Test P7B certificate",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/p7/cert2.p7b"},
			ExpectedResults: []CertificateInfo{{Name: "TestCert", Subject: "CN=cert2", Epoch: 1723889513, Type: "p7"}},
			ExpectedError:   "",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractP7CertificatesInfo); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test P7B certificate which is broken", func(t *testing.T) {
		tc := testCase{
			Name:            "Test P7B certificate which is broken",
			Cert:            Certificate{Name: "TestCert", Path: "../../tests/certs/p7/broken.p7b"},
			ExpectedResults: []CertificateInfo{},
			ExpectedError:   "Failed to decode any certificate in 'TestCert'",
		}
		if err := runExtractCertificateUnitTest(tc, t, ExtractP7CertificatesInfo); err != nil {
			t.Error(err)
		}
	})
}
