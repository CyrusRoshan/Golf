package ball

import (
	"time"

	"github.com/CyrusRoshan/Golf/physics"
	"github.com/CyrusRoshan/Golf/sectors"
	"github.com/CyrusRoshan/Golf/sprites"
	"github.com/faiface/pixel"
)

type Ball struct {
	physics physics.MotionFunctions

	sprite *pixel.Sprite
	bounds pixel.Rect

	initialTime time.Time

	initialVx float64
	initialVy float64

	initialX float64
	initialY float64
}

func NewBall(x float64, y float64) Ball {
	sprite, bounds := sprites.LoadSprite("img/ball.png")

	return Ball{
		initialX: x,
		initialY: y,

		sprite: sprite,
		bounds: bounds,

		physics: physics.MotionFunctions{
			X: physics.DistanceFunction(0, 0),
			Y: physics.DistanceFunction(0, physics.G),

			Vx: physics.VelocityFunction(0),
			Vy: physics.VelocityFunction(physics.G),
		},
	}
}

func (b *Ball) Draw(time float64, t pixel.Target) {
	movedVec := pixel.V(b.X(time), b.Y(time))
	b.sprite.Draw(t, pixel.IM.Moved(movedVec))
}

func (b *Ball) UpdateTrajectoryAtPoint(vx, vy, time float64) {
	x := b.physics.X(time)
	y := b.physics.Y(time)

	b.initialVx = vx
	b.initialVy = vy
	b.initialX = x
	b.initialY = y

	b.physics = physics.MotionFunctions{
		X: physics.DistanceFunction(vx, 0),
		Y: physics.DistanceFunction(vy, physics.G),

		Vx: physics.VelocityFunction(0),
		Vy: physics.VelocityFunction(physics.G),
	}
}

func (b *Ball) Vx(time float64) float64 {
	dVx := b.physics.Vx(time)
	return b.initialVx + dVx
}

func (b *Ball) Vy(time float64) float64 {
	dVy := b.physics.Vy(time)
	return b.initialVy + dVy
}

func (b *Ball) X(time float64) float64 {
	dX := b.physics.X(time)
	return b.initialX + dX
}

func (b *Ball) Y(time float64) float64 {
	dY := b.physics.Y(time)
	return b.initialY + dY
}

func (b *Ball) CollidesWith(f sectors.Func, timePrecision float64) (doesCollide bool, dt float64) {
	var ballX float64
	for dt = 0; ballX <= f.Range.EndX; dt += 1 {
		ballX = b.X(dt)
		ballY := b.Y(dt)

		pathY := f.F(ballX)

		if ballY < pathY {
			doesCollide = true

			minDt := dt - 1
			for ; dt > minDt; dt -= timePrecision {
				ballX := b.X(dt)
				ballY := b.Y(dt)

				pathY := f.F(ballX)

				if ballY >= pathY {
					return true, dt
				}
			}
		} else if ballY == pathY {
			return true, dt
		}
	}

	return false, 0
}
