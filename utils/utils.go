package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
