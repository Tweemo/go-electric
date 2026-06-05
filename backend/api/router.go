package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tweemo/go-electric/config"
	"github.com/tweemo/go-electric/rates"
)

// NewRouter builds the HTTP engine with middleware and registers all routes.
func NewRouter(cfg config.Config, r *rates.Rates) *gin.Engine {
	s := NewServer(cfg, r)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), corsMiddleware(cfg))
	engine.MaxMultipartMemory = cfg.MaxUploadBytes

	// Register routes here; new endpoints are one line each.
	engine.GET("/health", s.Health)
	engine.POST("/costs", s.Costs)

	return engine
}
