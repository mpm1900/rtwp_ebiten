package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type CollisionData struct {
	Size   dmath.Vec2
	Offset dmath.Vec2
}

type CollisionMoveResult struct {
	Moved    bool
	Collided bool
}

var Collision = donburi.NewComponentType[CollisionData]()
var CollisionQuery = donburi.NewQuery(filter.Contains(Collision, transform.Transform))

func WithCollision(entry *donburi.Entry, size dmath.Vec2) {
	WithCollisionOffset(entry, size, dmath.Vec2{})
}

func WithCollisionOffset(entry *donburi.Entry, size dmath.Vec2, offset dmath.Vec2) {
	entry.AddComponent(Collision)
	Collision.SetValue(entry, CollisionData{
		Size:   size,
		Offset: offset,
	})
}

func DetectCollisions(world donburi.World, yield func(a, b *donburi.Entry)) {
	colliders := []*donburi.Entry{}
	for entry := range CollisionQuery.Iter(world) {
		colliders = append(colliders, entry)
	}

	for i, a := range colliders {
		aBounds, ok := CollisionBounds(a)
		if !ok {
			continue
		}

		for _, b := range colliders[i+1:] {
			bBounds, ok := CollisionBounds(b)
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
	if !entry.HasComponent(transform.Transform) {
		return CollisionMoveResult{}
	}

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
	bounds, ok := CollisionBoundsAt(entry, position)
	if !ok {
		return nil, false
	}

	for other := range CollisionQuery.Iter(world) {
		if other.Entity() == entry.Entity() {
			continue
		}

		otherBounds, ok := CollisionBounds(other)
		if ok && bounds.Overlaps(otherBounds) {
			return other, true
		}
	}

	return nil, false
}

func CollisionBounds(entry *donburi.Entry) (image.Rectangle, bool) {
	if !entry.HasComponent(Collision) || !entry.HasComponent(transform.Transform) {
		return image.Rectangle{}, false
	}

	trans := transform.Transform.Get(entry)
	return CollisionBoundsAt(entry, trans.LocalPosition)
}

func CollisionCenter(entry *donburi.Entry) (dmath.Vec2, bool) {
	if !entry.HasComponent(Collision) || !entry.HasComponent(transform.Transform) {
		return dmath.Vec2{}, false
	}

	trans := transform.Transform.Get(entry)
	return CollisionCenterAt(entry, trans.LocalPosition)
}

func CollisionBoundsAt(entry *donburi.Entry, position dmath.Vec2) (image.Rectangle, bool) {
	if !entry.HasComponent(Collision) {
		return image.Rectangle{}, false
	}

	collision := Collision.Get(entry)
	if collision.Size.X <= 0 || collision.Size.Y <= 0 {
		return image.Rectangle{}, false
	}

	min := position.Add(collision.Offset)
	max := min.Add(collision.Size)

	return util.ToRect(min, max), true
}

func CollisionCenterAt(entry *donburi.Entry, position dmath.Vec2) (dmath.Vec2, bool) {
	if !entry.HasComponent(Collision) {
		return dmath.Vec2{}, false
	}

	collision := Collision.Get(entry)
	if collision.Size.X <= 0 || collision.Size.Y <= 0 {
		return dmath.Vec2{}, false
	}

	return position.Add(collision.Offset).Add(collision.Size.DivScalar(2)), true
}
