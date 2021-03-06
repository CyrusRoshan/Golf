package mathutils

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func PythagoreanC(a, b float64) float64 {
	return math.Sqrt(math.Pow(a, 2) + math.Pow(b, 2))
}

func Constrain(value, min, max float64) float64 {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

func ConstrainLineSlope(slope, intercept, width, minHeight, maxHeight float64) float64 {
	if slope*width+intercept > maxHeight {
		slope = (maxHeight - intercept) / width
	} else if slope*width+intercept < minHeight {
		slope = (minHeight - intercept) / width
	}

	return slope
}
