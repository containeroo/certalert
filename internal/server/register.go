package server

import "net/http"

// Handler is a struct that contains a route and a handler function.
type Handler struct {
	Path        string                                   // Path is the URL route for the handler.
	Description string                                   // Description provides a brief explanation of the handler's purpose.
	Handler     func(http.ResponseWriter, *http.Request) // Handler is the function that handles HTTP requests for the route.
	Methods     []string                                 // Methods is a list of HTTP methods supported by the handler.
}

// Handlers is a slice of Handler objects.
var Handlers = []Handler{}

// Register registers a handler function for a specific HTTP route.
//
// This function is used to register custom HTTP handlers along with their associated
// routes, descriptions, and supported HTTP methods. The registered handlers are later
// added to the router used by the HTTP server.
//
// Parameters:
//   - path: string
//     The URL route for the handler.
//   - h: func(http.ResponseWriter, *http.Request)
//     The function that handles HTTP requests for the route.
//   - description: string
//     A brief explanation of the handler's purpose.
//   - needsToken: bool
//     Indicates if the handler requires a valid API token.
//   - methods: ...string
//     Optional list of HTTP methods supported by the handler.
func Register(path, description string, handler func(http.ResponseWriter, *http.Request), methods ...string) {
	Handlers = append(Handlers,
		Handler{
			Path:        path,
			Description: description,
			Handler:     handler,
			Methods:     methods,
		})
}
