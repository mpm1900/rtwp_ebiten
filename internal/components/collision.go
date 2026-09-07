package components

import (
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var Collision = donburi.NewTag("Collision")
var CollisionQuery = donburi.NewQuery(filter.Contains(Collision, transform.Transform))

func WithCollision(entry *donburi.Entry) {
	entry.AddComponent(Collision)
}

func CollidesAt(world donburi.World, entry *donburi.Entry, position math.Vec2) (*donburi.Entry, bool) {
	bounds, ok := RectAt(entry, position)
	if !ok {
		return nil, false
	}

	var colliding *donburi.Entry
	collisionSpatial.eachCandidate(bounds, entry.Entity(), func(other *donburi.Entry) bool {
		other_bounds, ok := Rect(other)
		if ok && bounds.Overlaps(other_bounds) {
			colliding = other
			return false
		}

		return true
	})
	if colliding != nil {
		return colliding, true
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

func FirstColliderAtPoint(world donburi.World, point math.Vec2) (*donburi.Entry, bool) {
	pt := util.ToPoint(point)
	key := [2]int{cellCoord(pt.X), cellCoord(pt.Y)}

	for _, entry := range collisionSpatial.cells[key] {
		bounds, ok := Rect(entry)
		if !ok {
			continue
		}

		if pt.In(bounds) {
			return entry, true
		}
	}

	return nil, false
}
