package sectors

import (
	"image/color"
	"math"

	"github.com/CyrusRoshan/Golf/utils"
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

func GenerateSector(width float64, maxHeight float64, maxSegments int, color color.Color) (s Sector) {
	s.Width = width
	s.Color = color

	currentLength := 0.0
	currentHeight := 50.0

	segmentWidths := generateSegmentWidths(maxSegments, width)
	for _, sectorWidth := range segmentWidths {
		newFunc := RandLineFuncConstrained(-1.5, 1.5, MIN_STRUCTURE_HEIGHT, maxHeight, currentHeight, sectorWidth)
		newFunc.Width = sectorWidth

		currentLength += sectorWidth
		currentHeight = newFunc.F(sectorWidth)
		s.Funcs = append(s.Funcs, newFunc)
	}

	s.Polygons = triangulateSector(s)
	return s
}

func generateSegmentWidths(maxSegments int, maxWidth float64) (widths []float64) {
	for i := 0; i < maxSegments-1; i++ {
		var width float64

		if maxWidth > MIN_SEGMENT_WIDTH {
			width = utils.RandFloat(MIN_SEGMENT_WIDTH, math.Min(maxWidth, MAX_SEGMENT_WIDTH))
		} else if maxWidth > 0 {
			width = maxWidth
		} else {
			break
		}

		widths = append(widths, width)
		maxWidth -= width
	}

	return widths
}

func triangulateSector(s Sector) (polygons []*imdraw.IMDraw) {
	position := 0.0
	for _, segment := range s.Funcs {
		imd := imdraw.New(nil)
		imd.Color = s.Color

		bottomRight := pixel.V(position+segment.Width, 0)
		imd.Push(bottomRight)

		bottomLeft := pixel.V(position, 0)
		imd.Push(bottomLeft)

		topLeft := pixel.V(position, segment.F(0))
		imd.Push(topLeft)

		topRight := pixel.V(position+segment.Width, segment.F(segment.Width))
		imd.Push(topRight)

		position += segment.Width
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
