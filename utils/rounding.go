package utils

import (
	"fmt"
	"math"
)

// RoundFloat rounds a float64 to two decimal places
func RoundFloat(value float64, decimals int) (float64, error) {
	if decimals < 0 {
		return 0, fmt.Errorf("decimals must be non-negative")
	}

	// Check for NaN input
	if math.IsNaN(value) {
		return 0, fmt.Errorf("input value is NaN")
	}

	// Handle zero decimal places (return integer)
	if decimals == 0 {
		return math.Round(value), nil
	}

	pow := math.Pow(10, float64(decimals))
	rounded := math.Round(value*pow) / pow
	return rounded, nil
}
