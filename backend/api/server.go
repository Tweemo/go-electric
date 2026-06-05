// Package api wires the HTTP layer: middleware, routing, and request handlers.
// Add a new endpoint by writing a handler method on Server and registering it in
// router.go.
package api

import (
	"github.com/tweemo/go-electric/config"
	"github.com/tweemo/go-electric/rates"
)

// Server holds the dependencies shared by the HTTP handlers.
type Server struct {
	cfg   config.Config
	rates *rates.Rates
}

// NewServer constructs a Server with its dependencies.
func NewServer(cfg config.Config, r *rates.Rates) *Server {
	return &Server{cfg: cfg, rates: r}
}
