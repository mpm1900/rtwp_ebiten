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
	entry.AddComponent(Collision)
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

func CollisionStopDistance(entry *donburi.Entry, stopDistance float64) float64 {
	if !entry.HasComponent(Collision) || !entry.HasComponent(transform.Transform) {
		return stopDistance
	}

	scale := transform.Transform.Get(entry).LocalScale
	if scale.X <= 0 || scale.Y <= 0 {
		return stopDistance
	}

	return max(stopDistance, scale.Magnitude())
}
