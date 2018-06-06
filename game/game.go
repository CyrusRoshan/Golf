package game

import (
	"github.com/CyrusRoshan/Golf/screen"
	"github.com/CyrusRoshan/Golf/sectors"
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	NEARBY_SECTORS = 1
	SCALE          = 0.5
)

type Game struct {
	window *pixelgl.Window
	canvas *pixelgl.Canvas

	holeComplete *bool

	width  float64
	height float64

	currentSector *sectors.Sector
	nextSector    *sectors.Sector
}

func NewGame(win *pixelgl.Window) *Game {
	g := Game{
		window: win,
		canvas: screen.NewCanvas(screen.ScreenBounds()),
	}
	g.canvas.SetMatrix(pixel.IM.Moved(g.canvas.Bounds().Min).Scaled(pixel.ZV, SCALE))
	g.width = g.canvas.Bounds().W()
	g.height = g.canvas.Bounds().H()

	sectors.MIN_SEGMENT_WIDTH = g.width / 20
	sectors.MAX_SEGMENT_WIDTH = g.width / 4

	// set up sectors
	currentSector := sectors.GenerateSector(g.width, g.height-sectors.STRUCTURE_CEILING_GAP, 10, colornames.Orange)
	g.currentSector = &currentSector

	return &g
}

func (g *Game) Run() {
	for !g.window.Closed() {
		screen.LimitFPS(30, func() {
			screen.ScaleWindowToCanvas(g.window, g.canvas)
			g.window.Clear(colornames.Red)
			g.canvas.Clear(colornames.Black)

			g.DrawFrames()

			g.canvas.Draw(g.window, pixel.IM)
			g.window.Update()
		})
	}
}

// TODO: Draw transition from current to next sector
func (g *Game) DrawTransitions() {
	// drawSector := g.sectors.GetSector(0, 0)
	// moveVec := g.screenArea.Center().Sub(g.sectors.GetCenterOfSector(0, 0))
	// sectorMatrix := pixel.IM.Moved(moveVec)
}

func (g *Game) DrawFrames() {
	g.currentSector.Draw(g.canvas)
	// for _, object := range drawSector.Objects {
	// 	matrix := object.Physics.GetMatrix().Chained(sectorMatrix)
	// 	object.Sprite.DrawColorMask(g.canvas, matrix, object.Color)
	// }

	// g.canvas.Draw(g.window, pixel.IM.Moved(g.canvas.Bounds().Center()))
}
