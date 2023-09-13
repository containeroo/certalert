package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func init() {
	Register("/certificates", Certificates, "GET")
}

// Certificates is the handler for the /certificates route
// It fetches all the certificates and displays them in a tabular format
func Certificates(w http.ResponseWriter, r *http.Request) {
	certificatesInfo, err := certificates.Process(config.App.Certs, config.App.FailOnError)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tplData := TemplateData{
		CertInfos: certificatesInfo,
		CSS:       CSS,
		JS:        JS,
	}
	tpl, err := renderTemplate(tplBase, tplCertificates, tplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(tpl)); err != nil {
		log.Errorf("Unable to write response: %s", err)
	}
}
