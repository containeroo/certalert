package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"net/http"
)

func renderCertificates(certInfo []certificates.CertificateInfo) string {
	data := TemplateData{
		CertInfos: certInfo,
		CSS:       CSS,
		JS:        JS,
	}
	return renderTemplate(tplBase, tplCertificates, data)
}

func Certificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	certificatesInfo, err := certificates.Process(config.App.Certs, config.FailOnError)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(renderCertificates(certificatesInfo)))
}
