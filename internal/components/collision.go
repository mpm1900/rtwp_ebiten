package components

import (
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

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
