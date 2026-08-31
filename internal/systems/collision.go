package systems

import (
	"math"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

type CollisionMoveResult struct {
	Moved    bool
	Collided bool
}

func MoveWithCollision(world donburi.World, entry *donburi.Entry, delta dmath.Vec2) CollisionMoveResult {
	trans := transform.Transform.Get(entry)
	startPos := trans.LocalPosition
	if !delta.IsZero() {
		trans.LocalRotation = dmath.ToDegrees(math.Atan2(delta.Y, delta.X))
	}

	nextPosition := components.ClampWorldPosition(startPos.Add(delta))
	_, colliding := components.CollidesAt(world, entry, nextPosition)
	if !colliding {
		trans.LocalPosition = nextPosition
		return CollisionMoveResult{
			Moved: !delta.IsZero(),
		}
	}

	result := CollisionMoveResult{Collided: true}
	currentPos := startPos

	nextPosition = components.ClampWorldPosition(currentPos.Add(dmath.NewVec2(delta.X, 0)))
	_, colliding = components.CollidesAt(world, entry, nextPosition)
	if delta.X != 0 && !colliding {
		currentPos = nextPosition
		result.Moved = true
	}

	nextPosition = components.ClampWorldPosition(currentPos.Add(dmath.NewVec2(0, delta.Y)))
	_, colliding = components.CollidesAt(world, entry, nextPosition)
	if delta.Y != 0 && !colliding {
		currentPos = nextPosition
		result.Moved = true
	}

	trans.LocalPosition = components.ClampWorldPosition(currentPos)
	return result
}
