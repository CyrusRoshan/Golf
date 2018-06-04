package screen

import (
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const CANVAS_SCALE = 5

func ScreenBounds() (width float64, height float64) {
	return pixelgl.PrimaryMonitor().Size()
}

func NewCanvas(width float64, height float64) (canvas *pixelgl.Canvas) {
	scaleRect := pixel.R(-width/CANVAS_SCALE, -height/CANVAS_SCALE, width/CANVAS_SCALE, height/CANVAS_SCALE)

	canvas = pixelgl.NewCanvas(scaleRect)
	canvas.SetMatrix(pixel.IM)

	return canvas
}

func LimitFPS(fps int, loopContents func()) {
	d := time.Second / time.Duration(fps)

	select {
	case <-time.After(d):
		loopContents()
	}
}

var (
	frames = 0
	second = time.Tick(time.Second)
)

func ScaleWindowToCanvas(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	winBounds := win.Bounds()
	canBounds := canvas.Bounds()

	wRatio := winBounds.W() / canBounds.W()
	hRatio := winBounds.H() / canBounds.H()

	smallestRatio := math.Min(wRatio, hRatio)

	// scale window from canvas
	scaleMatrix := pixel.IM.Scaled(pixel.ZV, smallestRatio)
	win.SetMatrix(scaleMatrix.Moved(win.Bounds().Center()))

	// scale window bounds
	if wRatio != hRatio {
		newBounds := canBounds.ResizedMin(pixel.Vec{
			X: canBounds.W() * smallestRatio,
			Y: canBounds.H() * smallestRatio,
		})

		win.SetBounds(newBounds)
	}
}
