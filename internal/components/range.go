package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var Range = donburi.NewComponentType[float64]()
var RangeQuery = donburi.NewQuery(filter.And(
	filter.Contains(Range, transform.Transform),
	filter.Not(filter.Contains(Delay))),
)

func WithRange(entry *donburi.Entry, r float64) {
	entry.AddComponent(Range)
	Range.SetValue(entry, r)
}

func EachActorsInRange(world donburi.World, entry *donburi.Entry, yield func(*donburi.Entry)) {
	if !entry.HasComponent(transform.Transform) {
		return
	}

	entry_center := Center(entry)
	entry_range := *Range.Get(entry)
	for actor := range Actor.Iter(world) {
		actor_center := Center(actor)
		distance := entry_center.Distance(actor_center)
		if entry_range >= distance {
			yield(actor)
		}
	}
}

func EachActorFromPoint(world donburi.World, position math.Vec2, r float64, yield func(*donburi.Entry)) {
	for actor := range Actor.Iter(world) {
		distance := position.Distance(Center(actor))
		if r >= distance {
			yield(actor)
		}
	}
}

func NearestActorFromEntity(world donburi.World, entry *donburi.Entry) (*donburi.Entry, bool) {
	var nearest *donburi.Entry
	var dist *float64
	position := Center(entry)
	for actor := range Actor.Iter(world) {
		if actor.Entity() == entry.Entity() {
			continue
		}
		d := position.Distance(Center(actor))
		if dist == nil || *dist > d {
			dist = &d
			nearest = actor
		}
	}

	return nearest, dist != nil
}

func NearestLivingActorFromEntity(world donburi.World, entry *donburi.Entry) (*donburi.Entry, bool) {
	var nearest *donburi.Entry
	var dist *float64
	position := Center(entry)
	for actor := range Actor.Iter(world) {
		if actor.Entity() == entry.Entity() || !IsAlive(actor) {
			continue
		}

		d := position.Distance(Center(actor))
		if dist == nil || *dist > d {
			dist = &d
			nearest = actor
		}
	}

	return nearest, dist != nil
}
