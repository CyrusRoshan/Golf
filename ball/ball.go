package ball

import (
	"github.com/CyrusRoshan/Golf/physics"
	"github.com/CyrusRoshan/Golf/sectors"
)

type Ball struct {
	Physics physics.MotionFunctions

	initialVx float64
	initialVy float64

	initialX float64
	initialY float64
}

func (b *Ball) TeleportOver(dx, dy float64) {
	b.initialX += dx
	b.initialY += dy
}

func (b *Ball) UpdateTrajectory(vx, vy float64) {
	b.Physics = physics.MotionFunctions{
		X: physics.DistanceFunction(vx, 0),
		Y: physics.DistanceFunction(vy, physics.G),

		Vx: physics.VelocityFunction(0),
		Vy: physics.VelocityFunction(physics.G),
	}
}

func (b *Ball) Vx(time float64) float64 {
	dVx := b.Physics.Vx(time)
	return b.initialVx + dVx
}

func (b *Ball) Vy(time float64) float64 {
	dVy := b.Physics.Vy(time)
	return b.initialVy + dVy
}

func (b *Ball) X(time float64) float64 {
	dX := b.Physics.X(time)
	return b.initialX + dX
}

func (b *Ball) Y(time float64) float64 {
	dY := b.Physics.Y(time)
	return b.initialY + dY
}

func (b *Ball) CollidesWith(f sectors.Func, timePrecision float64) (doesCollide bool, dt float64) {
	var ballX float64
	for dt = 0; ballX <= f.Range.EndX; dt += 1 {
		ballX = b.Physics.X(dt)
		ballY := b.Physics.Y(dt)

		pathY := f.F(ballX)

		if ballY < pathY {
			doesCollide = true

			minDt := dt - 1
			for ; dt > minDt; dt -= timePrecision {
				ballX := b.Physics.X(dt)
				ballY := b.Physics.Y(dt)

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
