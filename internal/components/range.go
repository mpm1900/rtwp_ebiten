package components

import (
	"github.com/yohamta/donburi"
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
