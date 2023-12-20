package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/server"
	"net/http"
)

func init() {
	server.Register("/certificates", "Fetches and displays all the certificates in a tabular format", Certificates, "GET")
}

// Certificates is an HTTP handler that processes certificate information and renders an HTML page.
//
// This function retrieves certificate information based on the configuration settings and renders
// an HTML page displaying details about the certificates. If an error occurs during the processing
// or rendering, it returns an HTTP 500 Internal Server Error response.
//
// Parameters:
//   - w: http.ResponseWriter
//     The HTTP response writer.
//   - r: *http.Request
//     The HTTP request.
func Certificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

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

	w.Write([]byte(tpl))
}
