package game

import (
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
	collisionImpossible bool
	finishedAnimating   bool

	golfBall  *ball.Ball
	ballPaths []*ballPath

	currentSector *sectors.Sector
	nextSector    *sectors.Sector
}

func NewGame(win *pixelgl.Window) *Game {
	canvas := screen.NewCanvas(screen.ScreenBounds())
	canvas.SetMatrix(pixel.IM.Moved(canvas.Bounds().Min).Scaled(pixel.ZV, DEBUG_SCALE))

	g := Game{
		window: win,
		canvas: canvas,

		lastHitTime:       time.Now(),
		finishedAnimating: true,

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
	golfBall := ball.NewBall(sectors.MIN_SEGMENT_WIDTH, sectors.MAX_STRUCTURE_HEIGHT)
	g.golfBall = &golfBall

	return &g
}

func (g *Game) Run() {
	for !g.window.Closed() {
		screen.LimitFPS(10, func() {
			screen.ShowDebugFPS("Golf", g.window)
			g.currentSector.Draw(g.canvas)
			g.UpdateScreen()

			if g.finishedAnimating {
				if g.collisionImpossible {
					golfBall := ball.NewBall(sectors.MIN_SEGMENT_WIDTH, sectors.MAX_STRUCTURE_HEIGHT)
					g.golfBall = &golfBall

					g.collisionImpossible = false
					g.waitingForHit = false // has to drop down to ground first
				} else if g.waitingForHit {
					g.GetInput(g.secondsSinceLastHit()) // todo: fix this
				}

				if !g.collisionImpossible && !g.waitingForHit {
					g.ballPaths, g.collisionImpossible, g.waitingForHit = g.CalcualtePaths(g.golfBall, g.currentSector)
					g.lastHitTime = time.Now()
					g.finishedAnimating = false
				}
			}

			dt := g.secondsSinceLastHit()
			g.finishedAnimating = g.DrawCurrentPath(dt)
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

func (g *Game) DrawCurrentPath(dt float64) (finishedDrawing bool) {
	var i int
	var currentPath *ballPath
	for i, currentPath = range g.ballPaths {
		if dt >= currentPath.StartTime &&
			dt <= currentPath.EndTime {
			break
		}
	}

	time := dt - currentPath.StartTime

	if i == len(g.ballPaths)-1 {
		if g.collisionImpossible {
			ballX := currentPath.Ball.X(time)
			ballY := currentPath.Ball.Y(time)

			if !g.canvas.Bounds().Contains(pixel.V(ballX, ballY)) {
				finishedDrawing = true
			}
		} else if g.waitingForHit { // waiting for a hit
			time = 0
			finishedDrawing = true
		}
	}

	currentPath.Ball.Draw(time, g.canvas)
	return finishedDrawing
}

type ballPath struct {
	Ball      *ball.Ball
	StartTime float64
	EndTime   float64
}

// TODO: Move to physics
func (g *Game) CalcualtePaths(golfBall *ball.Ball, currentSector *sectors.Sector) (paths []*ballPath, collisionImpossible, waitingForHit bool) {
	currentTime := 0.0
	currentBall := golfBall

	lastPath := &ballPath{
		Ball:      currentBall,
		StartTime: currentTime,
	}
	paths = append(paths, lastPath)

	for {
		collisionCalculated, collisionTime, collisionSegment := currentBall.FindCollision(&currentSector.Segments, 0, 3)
		currentTime += collisionTime

		if !collisionCalculated {
			// flies out of screen
			collisionImpossible = true
			break
		}

		vx := currentBall.Vx(collisionTime)
		vy := currentBall.Vy(collisionTime)
		slope := collisionSegment.Slope(currentBall.X(collisionTime))

		newRawVx, newRawVy := physics.CollisionReflectionAngle(vx, vy, slope)
		newVx := physics.COLLISION_COEFFICIENT * newRawVx
		newVy := physics.COLLISION_COEFFICIENT * newRawVy

		newVx -= math.Copysign(physics.KINETIC_FRICTION, newVx)

		var xPaused, yPaused bool
		if math.Abs(newVx) <= physics.KINETIC_FRICTION {
			newVx = 0
			xPaused = true
		}
		if math.Abs(newVy) <= physics.MIN_Y_VELOCITY {
			newVy = 0
			yPaused = true
		}

		newBall := ball.NewBallWithVelocity(currentBall.X(collisionTime), currentBall.Y(collisionTime), newVx, newVy)
		currentBall = &newBall
		lastPath.EndTime = currentTime

		currentPath := &ballPath{
			Ball:      &newBall,
			StartTime: currentTime,
		}
		paths = append(paths, currentPath)

		if xPaused && yPaused {
			waitingForHit = true
			currentPath.EndTime = currentTime
			break
		}

		lastPath = currentPath
	}

	return paths, collisionImpossible, waitingForHit
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
		g.collisionImpossible = false
	}
}
