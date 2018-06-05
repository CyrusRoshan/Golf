package ball

import (
	"github.com/CyrusRoshan/Golf/physics"
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
