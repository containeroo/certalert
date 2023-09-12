package server

import "net/http"

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
	w.Write([]byte(body))
}
