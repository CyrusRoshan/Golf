package sectors

import (
	"image/color"
	"math"

	"github.com/CyrusRoshan/Golf/mathutils"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

const (
	STRUCTURE_CEILING_GAP = 20.0
	MIN_STRUCTURE_HEIGHT  = 0.0
	SECTOR_LINE_WIDTH     = 5
)

var (
	MIN_SEGMENT_WIDTH float64
	MAX_SEGMENT_WIDTH float64
)

type Sector struct {
	Color    color.Color
	Width    float64
	Funcs    []Func
	Polygons []*imdraw.IMDraw
}

func GenerateSector(sectorWidth float64, maxHeight float64, maxSegments int, color color.Color) (s Sector) {
	s.Width = sectorWidth
	s.Color = color

	currentLength := 0.0
	currentHeight := 50.0

	segmentWidths := generateSegmentWidths(maxSegments, sectorWidth)
	for _, width := range segmentWidths {
		unconstrainedSlope := mathutils.RandFloat(-1.5, 1.5)

		var slope float64
		if mathutils.RandFloat(0, 1) < 0.3 {
			slope = 0
		} else {
			slope = mathutils.ConstrainLineSlope(unconstrainedSlope, currentHeight, width, MIN_STRUCTURE_HEIGHT, maxHeight)
		}

		newFunc := NewRangedLineFunc(slope, currentLength, currentLength+width, currentHeight)

		currentLength += width
		currentHeight = newFunc.F(width)
		s.Funcs = append(s.Funcs, newFunc)
	}

	s.Polygons = triangulateSector(s)
	return s
}

func generateSegmentWidths(maxSegments int, maxWidth float64) (widths []float64) {
	remainingLength := maxWidth

	for i := 0; i < maxSegments-1; i++ {
		var width float64

		if remainingLength > MIN_SEGMENT_WIDTH {
			width = mathutils.RandFloat(MIN_SEGMENT_WIDTH, math.Min(maxWidth, MAX_SEGMENT_WIDTH))
			remainingLength -= width

			if remainingLength <= MIN_SEGMENT_WIDTH {
				width += remainingLength
				remainingLength = 0
			}
		} else {
			break
		}

		widths = append(widths, width)
	}

	return widths
}

func triangulateSector(s Sector) (polygons []*imdraw.IMDraw) {
	for _, segment := range s.Funcs {
		imd := imdraw.New(nil)
		imd.Color = s.Color

		imd.Push(pixel.V(segment.Range.StartX, segment.Range.StartY))
		imd.Push(pixel.V(segment.Range.EndX, segment.Range.EndY))
		imd.Push(pixel.V(segment.Range.EndX, 0))
		imd.Push(pixel.V(segment.Range.StartX, 0))

		imd.Polygon(SECTOR_LINE_WIDTH)

		polygons = append(polygons, imd)
	}

	return polygons
}

func (s *Sector) Draw(target pixel.Target) {
	for _, polygon := range s.Polygons {
		polygon.Draw(target)
	}
}
