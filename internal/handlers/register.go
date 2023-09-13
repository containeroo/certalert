package handlers

import "net/http"

// Handler is a struct that contains a route and a handler function
type Handler struct {
	Path        string
	Description string
	Handler     func(http.ResponseWriter, *http.Request)
	Methods     []string
}

// Handlers is a slice of Handler objects
var Handlers = []Handler{}

// RegisterHandler registers a handler with the Handlers slice
func Register(path string, description string, h func(http.ResponseWriter, *http.Request), methods []string) {
	Handlers = append(Handlers, Handler{Path: path, Description: description, Handler: h, Methods: methods})
}
