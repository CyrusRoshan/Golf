package physics

import (
	"math"
)

const (
	G                      = -9.8
	COLLISION_COEFFICIENT  = 0.7
	FRICTIONAL_COEFFICIENT = 0.6
)

type DeltaFunction func(float64) float64

type MotionFunctions struct {
	X  DeltaFunction
	Y  DeltaFunction
	Vx DeltaFunction
	Vy DeltaFunction
}

func DistanceFunction(v0 float64, accel float64) DeltaFunction {
	return func(time float64) float64 {
		return v0*time + (accel*math.Pow(time, 2))/2
	}
}

func VelocityFunction(accel float64) DeltaFunction {
	return func(time float64) float64 {
		return time * accel
	}
}

// not accurate, and combines both static and kinetic friction
func FrictionalForce(m, g float64) float64 {
	n := m * g
	return FRICTIONAL_COEFFICIENT * n
}

func CollisionReflectionAngle(vx, vy, slope float64) (newVx, newVy float64) {
	slopeAngle := math.Atan(slope)
	vAngle := math.Atan(vy / vx)

	newVangle := 2*slopeAngle - vAngle
	newUnscaledVx := math.Cos(newVangle)
	newUnscaledVy := math.Sin(newVangle)

	scaleMultiple := math.Sqrt(math.Pow(vx, 2) + math.Pow(vy, 2))
	newVx = newUnscaledVx * scaleMultiple
	newVy = newUnscaledVy * scaleMultiple

	return newVx, newVy
}
