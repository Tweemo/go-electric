package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tweemo/go-electric/cost_calculators"
	"github.com/tweemo/go-electric/utils"
)

// Costs accepts a multipart CSV upload (field "file") and returns the estimated
// cost of every plan. The upload is parsed in-memory; nothing is written to disk.
func (s *Server) Costs(c *gin.Context) {
	// The upload may contain usage figures the caller would rather not have cached
	// anywhere on the way back.
	c.Header("Cache-Control", "no-store")

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

	parsed, err := utils.ParseUsageData(src)
	if err != nil {
		// Parser errors are specific and safe to surface (e.g. "no valid usage
		// rows found", "could not find date/usage columns") so users can fix the file.
		slog.Error("parse usage data", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tell the client how much of their file was usable without changing the body shape.
	c.Header("X-Rows-Parsed", strconv.Itoa(parsed.RowsParsed))
	c.Header("X-Rows-Skipped", strconv.Itoa(parsed.RowsSkipped))

	records, err := utils.CalculateDayPower(parsed.Records)
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
