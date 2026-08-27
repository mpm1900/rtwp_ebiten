package components

import (
	"math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var Range = donburi.NewComponentType[float64]()
var RangeQuery = donburi.NewQuery(filter.Contains(Range, transform.Transform))

func WithRange(entry *donburi.Entry, r float64) {
	entry.AddComponent(Range)
	Range.SetValue(entry, r)
}

func EachActorsInRange(world donburi.World, entry *donburi.Entry, yield func(*donburi.Entry)) {
	if !entry.HasComponent(transform.Transform) {
		return
	}

	entry_transform := transform.Transform.Get(entry)
	entry_range := *Range.Get(entry)
	for actor := range ActorTag.Iter(world) {
		actor_transform := transform.Transform.Get(actor)
		distance := math.Hypot(
			entry_transform.LocalPosition.X-actor_transform.LocalPosition.X,
			entry_transform.LocalPosition.Y-actor_transform.LocalPosition.Y,
		)
		if entry_range >= distance {
			yield(actor)
		}
	}
}
