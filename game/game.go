package game

import (
	"time"

	"github.com/CyrusRoshan/Golf/ball"
	"github.com/CyrusRoshan/Golf/screen"
	"github.com/CyrusRoshan/Golf/sectors"
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	NEARBY_SECTORS = 1
	SCALE          = 0.5
	TIME_SCALE     = 3.0
)

type Game struct {
	window *pixelgl.Window
	canvas *pixelgl.Canvas

	holeComplete *bool

	width  float64
	height float64

	lastHitTime time.Time
	golfBall    ball.Ball

	currentSector *sectors.Sector
	nextSector    *sectors.Sector
}

func NewGame(win *pixelgl.Window) *Game {
	canvas := screen.NewCanvas(screen.ScreenBounds())
	canvas.SetMatrix(pixel.IM.Moved(canvas.Bounds().Min).Scaled(pixel.ZV, SCALE))

	g := Game{
		window: win,
		canvas: canvas,

		lastHitTime: time.Now(),
		golfBall:    ball.NewBall(10, 100),

		width:  canvas.Bounds().W(),
		height: canvas.Bounds().H(),
	}

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

	dt := float64(time.Now().Sub(g.lastHitTime))
	dtInSeconds := dt / float64(time.Second) * TIME_SCALE

	g.golfBall.Draw(dtInSeconds, g.canvas)
}
