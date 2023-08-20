package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	certificatesInfo, err := certificates.Process(config.App.Certs, config.FailOnError)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(renderCertificateInfo(certificatesInfo)))
}
