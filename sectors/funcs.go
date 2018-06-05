package sectors

import (
	"fmt"
	"math"

	"github.com/CyrusRoshan/Golf/utils"
)

type Func struct {
	Width float64
	F     func(float64) float64
	Df    func(float64) float64
}

func CalculateReflectionAngle(vx float64, vy float64, slope float64) (newVx float64, newVy float64) {
	slopeAngle := math.Atan(slope)
	vAngle := math.Atan(vy / vx)

	newVangle := (slopeAngle + 180 - vAngle) + slopeAngle
	newUnscaledVx := math.Cos(newVangle)
	newUnscaledVy := math.Sin(newVangle)

	scaleMultiple := math.Sqrt(math.Pow(vx, 2) + math.Pow(vy, 2))
	newVx = newUnscaledVx * scaleMultiple
	newVy = newUnscaledVy * scaleMultiple

	return newVx, newVy
}

func RandLineFuncConstrained(minSlope float64, maxSlope float64, minHeight float64, maxHeight float64, intercept float64, width float64) Func {
	slope := utils.RandFloat(minSlope, maxSlope)

	fmt.Println("-----")
	fmt.Println(slope)
	if slope*width+intercept > maxHeight {
		slope = (maxHeight - intercept) / width
	} else if slope*width+intercept < minHeight {
		slope = (minHeight - intercept) / width
	}
	fmt.Println(slope)

	return NewLineFunc(slope, intercept)
}

func NewLineFunc(slope float64, intercept float64) Func {
	f := func(x float64) float64 {
		return slope*x + intercept
	}
	df := func(x float64) float64 {
		return slope
	}

	return Func{
		F:  f,
		Df: df,
	}
}

func NewSinFunc(waveLength float64, waveHeight float64, offset float64, intercept float64, onlyPositive bool) Func {
	fWithoutOffset := func(x float64) float64 {
		return math.Sin(math.Pi*(x/waveLength+offset)) * waveHeight
	}

	interceptDiff := intercept - fWithoutOffset(0)
	f := func(x float64) float64 {
		return fWithoutOffset(x) + interceptDiff
	}

	df := func(x float64) float64 {
		return math.Sin(math.Pi*(x/waveLength+offset)) * waveHeight * math.Pi / waveLength
	}

	return Func{
		F:  f,
		Df: df,
	}
}
