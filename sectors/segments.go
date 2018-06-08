package sectors

import (
	"math"

	"github.com/faiface/pixel"
)

type Range struct {
	Start pixel.Vec

	End pixel.Vec
}

type Segment struct {
	Range Range

	Y     func(float64) float64 // TODO: rename to Y
	Slope func(float64) float64 // TODO: rename to Slope
}

func NewRangedLineSegment(slope, startX, endX, startY float64) Segment {
	yf := func(x float64) float64 {
		return slope*(x-startX) + startY
	}
	slopef := func(x float64) float64 {
		return slope
	}

	return Segment{
		Y:     yf,
		Slope: slopef,

		Range: Range{
			Start: pixel.V(startX, yf(startX)),
			End:   pixel.V(endX, yf(endX)),
		},
	}
}

func newLineSegment(slope, yIntercept float64) Segment {
	yf := func(x float64) float64 {
		return slope*x + yIntercept
	}
	slopef := func(x float64) float64 {
		return slope
	}

	return Segment{
		Y:     yf,
		Slope: slopef,
	}
}

func NewSinSegment(waveLength float64, waveHeight float64, offset float64, intercept float64, onlyPositive bool) Segment {
	fWithoutOffset := func(x float64) float64 {
		return math.Sin(math.Pi*(x/waveLength+offset)) * waveHeight
	}

	interceptDiff := intercept - fWithoutOffset(0)
	yf := func(x float64) float64 {
		return fWithoutOffset(x) + interceptDiff
	}

	slopef := func(x float64) float64 {
		return math.Sin(math.Pi*(x/waveLength+offset)) * waveHeight * math.Pi / waveLength
	}

	return Segment{
		Y:     yf,
		Slope: slopef,
	}
}
