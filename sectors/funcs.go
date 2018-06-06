package sectors

import (
	"math"
)

type Range struct {
	StartX float64
	EndX   float64

	StartY float64
	EndY   float64
}

type Func struct {
	Range Range

	F  func(float64) float64
	Df func(float64) float64
}

func NewRangedLineFunc(slope, startX, endX, startY float64) Func {
	f := func(x float64) float64 {
		return slope*(x-startX) + startY
	}
	df := func(x float64) float64 {
		return slope
	}

	return Func{
		F:  f,
		Df: df,

		Range: Range{
			StartX: startX,
			EndX:   endX,

			StartY: f(startX),
			EndY:   f(endX),
		},
	}
}

func NewLineFunc(slope, yIntercept float64) Func {
	f := func(x float64) float64 {
		return slope*x + yIntercept
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
