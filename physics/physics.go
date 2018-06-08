package physics

import (
	"math"

	"github.com/CyrusRoshan/Golf/mathutils"
)

const (
	G                     = -9.8
	COLLISION_COEFFICIENT = 0.7
	KINETIC_FRICTION      = 0.001
	MIN_Y_VELOCITY        = 1
)

type DeltaFunction func(float64) float64
type DoubleDeltaFunction func(float64) (float64, float64)

type DeltaFunctions struct {
	DX DeltaFunction
	DY DeltaFunction

	DVx DeltaFunction
	DVy DeltaFunction

	DTdx DoubleDeltaFunction
	DTdy DoubleDeltaFunction
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

func TimeFunction(v0 float64, accel float64) DoubleDeltaFunction {
	if accel == 0 {
		return func(position float64) (float64, float64) {
			time := position / v0
			return time, time
		}
	}

	return func(position float64) (float64, float64) {
		rightQuadraticHalf := math.Sqrt(math.Pow(v0, 2) + 2*accel*position)
		t1 := (-v0 - rightQuadraticHalf) / accel
		t2 := (-v0 + rightQuadraticHalf) / accel
		return t1, t2
	}
}

func CollisionReflectionAngle(vx, vy, slope float64) (newVx, newVy float64) {
	if slope == 0 {
		return vx, vy * -1
	}

	if slope == math.Inf(1) || slope == math.Inf(-1) {
		return vx * -1, vy
	}

	slopeAngle := math.Atan(slope)
	vAngle := math.Atan(vy / vx)

	newVangle := 2*slopeAngle - vAngle
	newUnscaledVx := math.Cos(newVangle)
	newUnscaledVy := math.Sin(newVangle)

	scaleMultiple := mathutils.PythagoreanC(vx, vy)
	newVx = newUnscaledVx * scaleMultiple
	newVy = newUnscaledVy * scaleMultiple

	return newVx, newVy
}
