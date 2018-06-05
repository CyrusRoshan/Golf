package physics

import "math"

const (
	G                          = -9.8
	AssumedFrictionCoefficient = 0.6
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
	return AssumedFrictionCoefficient * n
}
