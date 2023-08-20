package handlers

import (
	"net/http"
)

type Endpoint struct {
	Path        string   `json:"path"`
	Methods     []string `json:"method"`
	Description string   `json:"description"`
}

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

func renderEndpoints(endpoints []Endpoint) string {
	data := TemplateData{
		Endpoints: endpoints,
		CSS:       CSS,
	}
	return renderTemplate(tplBase, tplEndpoints, data)
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(renderEndpoints(Endpoints)))
}
