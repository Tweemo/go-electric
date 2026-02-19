package utils

import (
	"log"
	"os"
	"strconv"
)

// MustFloat64Env returns the float64 value of an environment variable, or panics if it's not set
func MustFloat64Env(key string) float64 {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("env %s is required and must be set", key)
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("env %s must be a valid number, got %q: %v", key, s, err)
	}
	return v
}
