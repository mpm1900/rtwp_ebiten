package components

import (
	"image"
	stdmath "math"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type CollisionData struct {
	Size   dmath.Vec2
	Offset dmath.Vec2
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

func CollisionBounds(entry *donburi.Entry) (image.Rectangle, bool) {
	if !entry.HasComponent(Collision) || !entry.HasComponent(transform.Transform) {
		return image.Rectangle{}, false
	}

	collision := Collision.Get(entry)
	if collision.Size.X <= 0 || collision.Size.Y <= 0 {
		return image.Rectangle{}, false
	}

	trans := transform.Transform.Get(entry)
	minX := trans.LocalPosition.X + collision.Offset.X
	minY := trans.LocalPosition.Y + collision.Offset.Y
	maxX := minX + collision.Size.X
	maxY := minY + collision.Size.Y

	return image.Rect(
		int(stdmath.Floor(minX)),
		int(stdmath.Floor(minY)),
		int(stdmath.Ceil(maxX)),
		int(stdmath.Ceil(maxY)),
	), true
}
