// Package config loads server configuration from the environment, applying the
// project's defaults. It is the single place env vars are read.
package config

import (
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// MaxUploadBytes is the largest multipart upload the server accepts (10 MiB).
const MaxUploadBytes int64 = 10 << 20

// Config holds the resolved server settings.
type Config struct {
	Host           string
	Port           string
	Env            string
	GinMode        string
	CORSOrigins    []string
	MaxUploadBytes int64
}

// Load reads configuration from the environment. Outside production it first
// loads a local .env file (best effort).
func Load() Config {
	env := os.Getenv("ENV")
	if env != "production" {
		_ = godotenv.Load()
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Host:           host,
		Port:           port,
		Env:            env,
		GinMode:        os.Getenv("GIN_MODE"),
		CORSOrigins:    corsOrigins(),
		MaxUploadBytes: MaxUploadBytes,
	}
}

// Addr returns the host:port the server should listen on.
func (c Config) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func corsOrigins() []string {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if strings.TrimSpace(raw) == "" {
		// Dev default — set explicitly in production.
		return []string{"http://localhost:3001", "http://127.0.0.1:3001"}
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			origins = append(origins, p)
		}
	}
	return origins
}
