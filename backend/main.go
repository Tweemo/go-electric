package main

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tweemo/go-electric/cost_calculators"
	"github.com/tweemo/go-electric/utils"
)

func corsMiddleware() gin.HandlerFunc {
	origins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	if len(origins) == 1 && origins[0] == "" {
		// Dev default — set explicitly in production
		origins = []string{"http://localhost:3001", "http://127.0.0.1:3001"}
	}
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // see below
		MaxAge:           12 * time.Hour,
	})
}

func main() {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	router := gin.Default()
	router.MaxMultipartMemory = 10 << 20 // 10 MiB
	router.Use(corsMiddleware())

	router.POST("/costs", Costs)

	// For some reason none of these envs are loading
	port := os.Getenv("PORT")
	if port == "" {
		// port = "3000"
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	addr := net.JoinHostPort(host, port)

	router.Run(addr)
}

func Costs(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	if err := os.MkdirAll("data", 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating data directory"})
		return
	}

	dataFilePath := filepath.Join("data", "data.csv")
	dst, err := os.Create(dataFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating the file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving the file"})
		return
	}

	data, err := utils.GetUsageData(dataFilePath)
	if err != nil {
		slog.Error("usage data", "err", err, "request_id", "todo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or unreadable CSV"})
		return
	}

	sortedRecords, _ := utils.CalculateDayPower(data)
	response := cost_calculators.AllPrices(sortedRecords)

	c.JSON(http.StatusOK, response)
}
