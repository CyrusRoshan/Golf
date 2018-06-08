package game

import (
	"fmt"
	"math"
	"time"

	"github.com/CyrusRoshan/Golf/ball"
	"github.com/CyrusRoshan/Golf/physics"
	"github.com/CyrusRoshan/Golf/screen"
	"github.com/CyrusRoshan/Golf/sectors"
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	DEBUG_SCALE = 0.7
	TIME_SCALE  = 15.0
)

type Game struct {
	window *pixelgl.Window
	canvas *pixelgl.Canvas

	holeComplete *bool

	width  float64
	height float64

	lastHitTime         time.Time
	waitingForHit       bool
	collisionCalculated bool
	collisionTime       float64
	collisionSegment    *sectors.Segment
	collisionImpossible bool

	golfBall ball.Ball

	currentSector *sectors.Sector
	nextSector    *sectors.Sector
}

func NewGame(win *pixelgl.Window) *Game {
	canvas := screen.NewCanvas(screen.ScreenBounds())
	canvas.SetMatrix(pixel.IM.Moved(canvas.Bounds().Min).Scaled(pixel.ZV, DEBUG_SCALE))

	g := Game{
		window: win,
		canvas: canvas,

		lastHitTime: time.Now(),

		width:  canvas.Bounds().W(),
		height: canvas.Bounds().H(),
	}

	sectors.MIN_SEGMENT_WIDTH = g.width / 20
	sectors.MAX_SEGMENT_WIDTH = g.width / 4
	sectors.MAX_STRUCTURE_HEIGHT = g.height - sectors.STRUCTURE_CEILING_GAP

	// set up sectors
	currentSector := sectors.GenerateSector(g.width, sectors.MAX_STRUCTURE_HEIGHT, 10, colornames.Orange)
	g.currentSector = &currentSector

	// set up ball
	g.golfBall = ball.NewBall(sectors.MIN_SEGMENT_WIDTH, sectors.MAX_STRUCTURE_HEIGHT)

	return &g
}

func (g *Game) Run() {
	for !g.window.Closed() {
		screen.LimitFPS(30, func() {
			screen.ShowDebugFPS("Golf", g.window)

			fmt.Println("A")
			g.GetInput(g.secondsSinceLastHit())
			fmt.Println("B")
			g.CalcualteDeltas(g.secondsSinceLastHit())
			fmt.Println("C")
			g.DrawFrames(g.secondsSinceLastHit())
			fmt.Println("D")

			g.UpdateScreen()
		})
	}
}

func (g *Game) UpdateScreen() {
	screen.ScaleWindowToCanvas(g.window, g.canvas)

	g.canvas.Draw(g.window, pixel.IM)
	g.window.Update()

	g.window.Clear(colornames.Red)
	g.canvas.Clear(colornames.Black)
}

// TODO: Draw transition from current to next sector
func (g *Game) DrawTransitions() {
	// drawSector := g.sectors.GetSector(0, 0)
	// moveVec := g.screenArea.Center().Sub(g.sectors.GetCenterOfSector(0, 0))
	// sectorMatrix := pixel.IM.Moved(moveVec)
}

func (g *Game) DrawFrames(dt float64) {
	g.golfBall.Draw(dt, g.canvas)
	g.currentSector.Draw(g.canvas)
}

func (g *Game) CalcualteDeltas(dt float64) {
	if g.collisionImpossible || g.waitingForHit {
		return
	}

	if !g.collisionCalculated {
		g.collisionCalculated, g.collisionTime, g.collisionSegment = g.golfBall.FindCollision(&g.currentSector.Segments, dt, 3)
		fmt.Println("CURRENT AND COLLISION TIMES", dt, g.collisionTime)

		if !g.collisionCalculated {
			fmt.Println("NO COLLISION!")
			g.collisionImpossible = true
		}
	}

	if g.collisionCalculated {
		if dt >= g.collisionTime {
			vx := g.golfBall.Vx(g.collisionTime)
			vy := g.golfBall.Vy(g.collisionTime)
			slope := g.collisionSegment.Slope(g.collisionTime)

			newRawVx, newRawVy := physics.CollisionReflectionAngle(vx, vy, slope)
			newVx := physics.COLLISION_COEFFICIENT * newRawVx
			newVy := physics.COLLISION_COEFFICIENT * newRawVy

			newVx -= math.Copysign(physics.KINETIC_FRICTION, newVx)
			fmt.Println("NEWV", newVx, newVy)

			var xPaused, yPaused bool
			if math.Abs(newVx) <= physics.KINETIC_FRICTION {
				newVx = 0
				xPaused = true
			}
			if math.Abs(newVy) <= physics.MIN_Y_VELOCITY {
				newVy = 0
				yPaused = true
			}
			if xPaused && yPaused {
				g.waitingForHit = true
			}

			fmt.Println("velocity", vx, vy)
			fmt.Println("slope", slope)
			fmt.Println("NEWVRAWS,", newRawVx, newRawVy)
			fmt.Println("NEWV", newVx, newVy)
			g.golfBall.UpdateTrajectoryAtTime(newVx, newVy, g.collisionTime)

			g.lastHitTime = time.Now()
			g.collisionCalculated = false
			g.collisionTime = 0
			g.collisionSegment = nil
		}
	}
}

func (g *Game) secondsSinceLastHit() float64 {
	dt := float64(time.Now().Sub(g.lastHitTime))
	return dt / float64(time.Second) * TIME_SCALE
}

func (g *Game) GetInput(dt float64) {
	vx, vy := g.golfBall.Vx(dt), g.golfBall.Vy(dt)

	changed := false

	if g.window.Pressed(pixelgl.KeyLeft) || g.window.Pressed(pixelgl.KeyA) {
		vx -= 10.0
		changed = true
	}
	if g.window.Pressed(pixelgl.KeyRight) || g.window.Pressed(pixelgl.KeyD) {
		vx += 10.0
		changed = true
	}

	if g.window.Pressed(pixelgl.KeyUp) || g.window.Pressed(pixelgl.KeyW) {
		vy += 10.0
		changed = true
	}
	if g.window.Pressed(pixelgl.KeyDown) || g.window.Pressed(pixelgl.KeyS) {
		vy -= 10.0
		changed = true
	}

	if changed {
		g.golfBall.UpdateTrajectoryAtTime(vx, vy, dt)
		g.lastHitTime = time.Now()
		g.collisionCalculated = false
		g.collisionImpossible = false
	}
}
