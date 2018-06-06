package ball

import (
	"time"

	"github.com/CyrusRoshan/Golf/physics"
	"github.com/CyrusRoshan/Golf/sectors"
	"github.com/CyrusRoshan/Golf/sprites"
	"github.com/faiface/pixel"
)

type Ball struct {
	physics physics.DeltaFunctions

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

		physics: physics.DeltaFunctions{
			DX: physics.DistanceFunction(0, 0),
			DY: physics.DistanceFunction(0, physics.G),

			DVx: physics.VelocityFunction(0),
			DVy: physics.VelocityFunction(physics.G),
		},
	}
}

func (b *Ball) Draw(time float64, t pixel.Target) {
	movedVec := pixel.V(b.X(time), b.Y(time))
	b.sprite.Draw(t, pixel.IM.Moved(movedVec))
}

func (b *Ball) UpdateTrajectoryAtTime(vx, vy, time float64) {
	b.initialVx = vx
	b.initialVy = vy
	b.initialX = b.X(time)
	b.initialY = b.Y(time)

	b.physics = physics.DeltaFunctions{
		DX: physics.DistanceFunction(vx, 0),
		DY: physics.DistanceFunction(vy, physics.G),

		DVx: physics.VelocityFunction(0),
		DVy: physics.VelocityFunction(physics.G),
	}
}

func (b *Ball) Vx(time float64) float64 {
	dVx := b.physics.DVx(time)
	return b.initialVx + dVx
}

func (b *Ball) Vy(time float64) float64 {
	dVy := b.physics.DVy(time)
	return b.initialVy + dVy
}

func (b *Ball) X(time float64) float64 {
	dX := b.physics.DX(time)
	return b.initialX + dX
}

func (b *Ball) Y(time float64) float64 {
	dY := b.physics.DY(time)
	return b.initialY + dY
}

// can be optimized by using binary search
func (b *Ball) CollidesWith(segment sectors.Segment, timePrecision float64) (doesCollide bool, dt float64) {
	if b.initialVx > 0 && b.initialX > segment.Range.EndX ||
		b.initialVx < 0 && b.initialX < segment.Range.StartX {
		return
	}

	ballWasInRange := false
	for dt = 0; true; dt += timePrecision {
		ballX := b.X(dt)
		ballY := b.Y(dt)

		if ballX >= segment.Range.StartX && ballX <= segment.Range.EndX {
			ballWasInRange = true
		} else {
			if ballWasInRange {
				break
			}
			continue
		}

		pathY := segment.F(ballX)

		if ballY < pathY {
			doesCollide = true

			minDt := dt - 1
			for ; dt > minDt; dt -= timePrecision / 2 {
				ballX := b.X(dt)
				ballY := b.Y(dt)

				pathY := segment.F(ballX)

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
