package components

import (
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type CollisionMoveResult struct {
	Moved    bool
	Collided bool
}

var Collision = donburi.NewTag("Collision")
var CollisionQuery = donburi.NewQuery(filter.Contains(Collision, transform.Transform))

func WithCollision(entry *donburi.Entry) {
	if !entry.HasComponent(Collision) {
		entry.AddComponent(Collision)
	}
}

func DetectCollisions(world donburi.World, yield func(a, b *donburi.Entry)) {
	colliders := []*donburi.Entry{}
	for entry := range CollisionQuery.Iter(world) {
		colliders = append(colliders, entry)
	}

	for i, a := range colliders {
		aBounds, ok := Rect(a)
		if !ok {
			continue
		}

		for _, b := range colliders[i+1:] {
			bBounds, ok := Rect(b)
			if !ok {
				continue
			}

			if aBounds.Overlaps(bBounds) {
				yield(a, b)
			}
		}
	}
}

func MoveWithCollision(world donburi.World, entry *donburi.Entry, delta dmath.Vec2) CollisionMoveResult {
	trans := transform.Transform.Get(entry)
	position := trans.LocalPosition
	nextPosition := position.Add(delta)
	_, colliding := CollidesAt(world, entry, nextPosition)
	if !colliding {
		trans.LocalPosition = nextPosition
		return CollisionMoveResult{
			Moved: !delta.IsZero(),
		}
	}

	result := CollisionMoveResult{Collided: true}
	nextPosition = position.Add(dmath.NewVec2(delta.X, 0))
	_, colliding = CollidesAt(world, entry, nextPosition)
	if delta.X != 0 && !colliding {
		position = nextPosition
		result.Moved = true
	}

	nextPosition = position.Add(dmath.NewVec2(0, delta.Y))
	_, colliding = CollidesAt(world, entry, nextPosition)
	if delta.Y != 0 && !colliding {
		position = nextPosition
		result.Moved = true
	}

	trans.LocalPosition = position
	return result
}

func CollidesAt(world donburi.World, entry *donburi.Entry, position dmath.Vec2) (*donburi.Entry, bool) {
	bounds, ok := RectAt(entry, position)
	if !ok {
		return nil, false
	}

	for other := range CollisionQuery.Iter(world) {
		if other.Entity() == entry.Entity() {
			continue
		}

		otherBounds, ok := Rect(other)
		if ok && bounds.Overlaps(otherBounds) {
			return other, true
		}
	}

	return nil, false
}
