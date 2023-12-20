package handlers

import (
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/server"
	"net/http"
)

func init() {
	server.Register("/healthz", "Returns the health of the application", Healthz, "GET", "POST")
}

// Healthz is an HTTP handler for the /healthz route.
//
// This handler returns the status of the application. It checks the health
// of the application by attempting to process the configured certificates. If
// the certificate processing encounters an error, it returns an HTTP 500
// Internal Server Error response with the error message. If the application is
// healthy, it returns an HTTP 200 OK response with the "ok" message.
//
// Parameters:
//   - w: http.ResponseWriter
//     The HTTP response writer.
//   - r: *http.Request
//     The HTTP request.
func Healthz(w http.ResponseWriter, r *http.Request) {
	if _, err := certificates.Process(config.App.Certs, true); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
