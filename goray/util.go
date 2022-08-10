package main

import (
	"math/rand"
)

func RandomRange(min, max float64) float64 {
	// Returns a random real in [min,max).
	return min + (max-min)*rand.Float64()
}