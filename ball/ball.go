package ball

import (
	"math"

	"github.com/CyrusRoshan/Golf/physics"
	"github.com/CyrusRoshan/Golf/sectors"
	"github.com/CyrusRoshan/Golf/sprites"
	"github.com/faiface/pixel"
)

type Ball struct {
	physics physics.DeltaFunctions

	sprite *pixel.Sprite
	bounds pixel.Rect

	initialVx float64
	initialVy float64

	initialX float64
	initialY float64
}

func NewBall(x, y float64) Ball {
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

			DTdx: physics.TimeFunction(0, 0),
			DTdy: physics.TimeFunction(0, physics.G),
		},
	}
}

func NewBallWithVelocity(x, y, vx, vy float64) Ball {
	sprite, bounds := sprites.LoadSprite("img/ball.png")

	return Ball{
		initialX: x,
		initialY: y,

		initialVx: vx,
		initialVy: vy,

		sprite: sprite,
		bounds: bounds,

		physics: physics.DeltaFunctions{
			DX: physics.DistanceFunction(vx, 0),
			DY: physics.DistanceFunction(vy, physics.G),

			DVx: physics.VelocityFunction(0),
			DVy: physics.VelocityFunction(physics.G),

			DTdx: physics.TimeFunction(vx, 0),
			DTdy: physics.TimeFunction(vy, physics.G),
		},
	}
}

func (b *Ball) Draw(time float64, t pixel.Target) {
	movedVec := pixel.V(b.X(time), b.Y(time))
	b.sprite.Draw(t, pixel.IM.Moved(movedVec))
}

func (b *Ball) UpdateTrajectoryAtTime(vx, vy, time float64) {
	b.initialX = b.X(time)
	b.initialY = b.Y(time)
	b.initialVx = vx
	b.initialVy = vy

	b.physics = physics.DeltaFunctions{
		DX: physics.DistanceFunction(vx, 0),
		DY: physics.DistanceFunction(vy, physics.G),

		DVx: physics.VelocityFunction(0),
		DVy: physics.VelocityFunction(physics.G),

		DTdx: physics.TimeFunction(vx, 0),
		DTdy: physics.TimeFunction(vy, physics.G),
	}
}

func (b *Ball) Vx(time float64) float64 {
	dvx := b.physics.DVx(time)
	return b.initialVx + dvx
}
func (b *Ball) Vy(time float64) float64 {
	dvy := b.physics.DVy(time)
	return b.initialVy + dvy
}

func (b *Ball) X(time float64) float64 {
	dx := b.physics.DX(time)
	return b.initialX + dx
}
func (b *Ball) Y(time float64) float64 {
	dy := b.physics.DY(time)
	return b.initialY + dy
}

func (b *Ball) TgivenX(x float64) float64 {
	dt, _ := b.physics.DTdx(x - b.initialX)
	return dt
}
func (b *Ball) TgivenY(y float64) (float64, float64) {
	dt1, dt2 := b.physics.DTdy(y - b.initialY)
	return dt1, dt2
}

func (b *Ball) FindCollision(segments *[]sectors.Segment, startTime, heightPrecision float64) (doesCollide bool, dt float64, collisionSegment *sectors.Segment) {
	searchDirection := 1
	if b.initialVx < 0 {
		searchDirection = -1
	}

	// find the sector the ball is currently hovering over
	currentBallX := b.X(startTime)
	startSector := 0
	overSector := false
	for i, segment := range *segments {
		if segment.Range.Start.X <= currentBallX && currentBallX <= segment.Range.End.X {
			startSector = i
			overSector = true
		}
	}

	// for ball dropping straight down
	if b.initialVx == 0 {
		if !overSector {
			return false, 0, nil
		}

		y := ((*segments)[startSector]).Y(b.initialX)
		dt = math.Max(b.TgivenY(y + heightPrecision))

		if dt < 0 {
			return false, 0, nil
		}
		return true, dt, &((*segments)[startSector])
	}

	// check all sectors the ball is traveling towards, in order, for collisions
	for i := startSector; i >= 0 && i < len(*segments); i += searchDirection {
		segment := (*segments)[i]

		var startPosition, endPosition pixel.Vec
		if searchDirection > 0 {
			startPosition = segment.Range.Start
			endPosition = segment.Range.End
		} else {
			startPosition = segment.Range.End
			endPosition = segment.Range.Start
		}

		// find range of time that the ball can reasonbly be colliding with the segment in
		minTime := math.Max(startTime, b.TgivenX(startPosition.X))
		maxTime := b.TgivenX(endPosition.X)

		if b.Y(maxTime) > endPosition.Y && // ball is flying over this segment
			b.Y(minTime) > segment.Y(b.X(minTime)) {
			continue
		}

		for {
			midpointTime := (minTime + maxTime) / 2
			ballY := b.Y(midpointTime)
			ballX := b.X(midpointTime)
			segmentY := segment.Y(ballX)

			distanceAboveSegment := ballY - segmentY
			if distanceAboveSegment > 0 &&
				distanceAboveSegment < heightPrecision {
				dt = midpointTime
				doesCollide = true
				break
			}

			if distanceAboveSegment > 0 {
				minTime = midpointTime
			} else {
				maxTime = midpointTime
			}
		}

		if doesCollide {
			return true, dt, &((*segments)[i])
		}
	}

	return false, 0, nil
}
