package components

import (
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type CollisionData struct {
	OnCollide func(donburi.World, *donburi.Entry, donburi.Entity)
}

var Collision = donburi.NewComponentType[CollisionData]()
var CollisionQuery = donburi.NewQuery(filter.Contains(Collision, transform.Transform))

func WithCollision(entry *donburi.Entry, data CollisionData) {
	entry.AddComponent(Collision)
	Collision.SetValue(entry, data)
}

func CollidesAt(world donburi.World, entry *donburi.Entry, position math.Vec2) (*donburi.Entry, bool) {
	entity := entry.Entity()
	if !world.Valid(entity) {
		return nil, false
	}

	entry = world.Entry(entity)
	bounds, ok := RectAt(entry, position)
	if !ok {
		return nil, false
	}

	var colliding *donburi.Entry
	collisionSpatial.eachCandidate(world, bounds, entity, func(other *donburi.Entry) bool {
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

	for _, entity := range collisionSpatial.cells[key] {
		if !world.Valid(entity) {
			continue
		}

		entry := world.Entry(entity)
		if !entry.HasComponent(transform.Transform) {
			continue
		}

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
