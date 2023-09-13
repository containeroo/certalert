package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func init() {
	Register("/", Home, "GET", "POST")
}

// Endpoint is a struct that represents an endpoint
type Endpoint struct {
	Path        string   `json:"path"`
	Methods     []string `json:"method"`
	Description string   `json:"description"`
}

// Endpoints is a list of all the endpoints
var Endpoints = []Endpoint{
	{
		Path:        "/",
		Methods:     []string{"GET", "POST"},
		Description: "Shows this page (the endpoints)",
	},
	{
		Path:        "/certificates",
		Methods:     []string{"GET", "POST"},
		Description: "Fetches and displays all the certificates in a tabular format",
	},
	{
		Path:        "/-/reload",
		Methods:     []string{"GET", "POST"},
		Description: "Reloads the configuration",
	},
	{
		Path:        "/config",
		Methods:     []string{"GET", "POST"},
		Description: "Provides the currently active configuration file. Plaintext passwords are redacted",
	},
	{
		Path:        "/metrics",
		Methods:     []string{"GET", "POST"},
		Description: "Delivers metrics for Prometheus to scrape",
	},
	{
		Path:        "/healthz",
		Methods:     []string{"GET", "POST"},
		Description: "Returns the health of the application",
	},
}

// Home is the handler for the / route
// It displays all the endpoints
func Home(w http.ResponseWriter, r *http.Request) {
	tplData := TemplateData{
		Endpoints: Endpoints,
		CSS:       CSS,
	}
	tpl, err := renderTemplate(tplBase, tplEndpoints, tplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(tpl)); err != nil {
		log.Errorf("Unable to write response: %s", err)
	}
}
