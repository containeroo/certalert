package server

import (
	"context"
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

var StaticFS embed.FS

// Run starts the HTTP server.
//
// Parameters:
//   - listenAddress: string
//     The address on which the server should listen (e.g., ":8080").
func Run(listenAddress string) {
	// Use a buffered channel for runChan to prevent signal drops
	runChan := make(chan os.Signal, 1)
	signal.Notify(runChan, os.Interrupt, syscall.SIGTERM)

	// Create a cancelable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the router and configure routes
	router := newRouter()

	// Create the HTTP server
	server := &http.Server{
		Addr:         listenAddress,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	// Run the server in a new goroutine
	go func() {
		log.Info().Msgf("Server is starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Server failed to start: %v", err)
		}
	}()

	// Wait for signals
	sig := <-runChan
	log.Info().Msgf("Received signal: %v. Shutting down gracefully...", sig)
	cancel() // Cancel context to trigger graceful shutdown
	// Force shutdown after 10 seconds
	time.AfterFunc(10*time.Second, func() {
		log.Fatal().Msg("Timed out waiting for server to shut down")
	})

	// Shutdown the server and wait for it to finish
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Msgf("Error during server shutdown: %v", err)
	}

	log.Info().Msg("Server gracefully shut down")
}

// newRouter generates the router used in the HTTP Server.
//
// Returns:
//   - *mux.Router
//     A configured instance of the Gorilla Mux router.
func newRouter() *mux.Router {
	router := mux.NewRouter()

	// Add the handlers to the router
	for _, h := range Handlers {
		router.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}

	// Custom 404 handler
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return router
}

// notFoundHandler handles 404 responses.
//
// Parameters:
//   - w: http.ResponseWriter
//     The HTTP response writer.
//   - r: *http.Request
//     The HTTP request being processed.
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
		log.Error().Msgf("Error writing 404 response: %v", err)
	}
}
