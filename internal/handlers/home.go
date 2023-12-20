package handlers

import (
	"certalert/internal/server"
	"net/http"
)

func init() {
	server.Register("/", "Shows this page (the endpoints)", Home, "GET", "POST")
}

// Home is an HTTP handler for the / route.
//
// This handler displays all the available endpoints in an HTML format. It sets
// the "Content-Type" header to "text/html" and renders the HTML template with
// information about the registered endpoints.
//
// Parameters:
//   - w: http.ResponseWriter
//     The HTTP response writer.
//   - r: *http.Request
//     The HTTP request.
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tplData := TemplateData{
		Endpoints: server.Handlers,
		CSS:       CSS,
	}
	tpl, err := renderTemplate(tplBase, tplEndpoints, tplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tpl))
}
