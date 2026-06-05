package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tweemo/go-electric/cost_calculators"
	"github.com/tweemo/go-electric/utils"
)

// Costs accepts a multipart CSV upload (field "file") and returns the estimated
// cost of every plan. The upload is parsed in-memory; nothing is written to disk.
func (s *Server) Costs(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file upload"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not open uploaded file"})
		return
	}
	defer src.Close()

	data, err := utils.ParseUsageData(src)
	if err != nil {
		slog.Error("parse usage data", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or unreadable CSV"})
		return
	}

	records, err := utils.CalculateDayPower(data)
	if err != nil {
		slog.Error("calculate day power", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process usage data"})
		return
	}

	response, err := cost_calculators.AllPrices(records, s.rates)
	if err != nil {
		slog.Error("calculate costs", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not calculate costs"})
		return
	}

	c.JSON(http.StatusOK, response)
}
