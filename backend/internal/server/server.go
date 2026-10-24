// Package server runs and shuts down the warehouse HTTP server.
package server

import (
	"context"
	"net/http"
	"time"
)

// Server owns the HTTP server and its shutdown lifecycle.
type Server struct {
	httpServer *http.Server
}

// Run starts serving HTTP requests on the given port.
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the running HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
