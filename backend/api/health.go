package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health is a simple liveness probe.
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
