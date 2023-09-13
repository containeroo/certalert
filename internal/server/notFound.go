package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// NotFoundHandler is a handler for 404 errors
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	body := `
		<html>
			<head>
				<title>404 Not Found</title>
			</head>
			<body>
				<h1>404 Not Found</h1>
				<p>The page you requested could not be found.</p>
			</body>
		</html>
	`
	if _, err := w.Write([]byte(body)); err != nil {
		log.Errorf("Unable to write response: %s", err)
	}
}
