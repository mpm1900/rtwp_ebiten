package systems

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

type CollisionMoveResult struct {
	Moved    bool
	Collided bool
}

func MoveWithCollision(world donburi.World, entry *donburi.Entry, delta math.Vec2) CollisionMoveResult {
	trans := transform.Transform.Get(entry)
	position := trans.LocalPosition
	nextPosition := position.Add(delta)
	_, colliding := components.CollidesAt(world, entry, nextPosition)
	if !colliding {
		trans.LocalPosition = nextPosition
		return CollisionMoveResult{
			Moved: !delta.IsZero(),
		}
	}

	result := CollisionMoveResult{Collided: true}
	nextPosition = position.Add(math.NewVec2(delta.X, 0))
	_, colliding = components.CollidesAt(world, entry, nextPosition)
	if delta.X != 0 && !colliding {
		position = nextPosition
		result.Moved = true
	}

	nextPosition = position.Add(math.NewVec2(0, delta.Y))
	_, colliding = components.CollidesAt(world, entry, nextPosition)
	if delta.Y != 0 && !colliding {
		position = nextPosition
		result.Moved = true
	}

	trans.LocalPosition = position
	return result
}
