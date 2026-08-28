package components

import (
	"image"
	"rtwp_ebitengine/internal/assets"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var ActorTag = donburi.NewTag("Actor")
var ActorQuery = donburi.NewQuery(filter.Contains(ActorTag, transform.Transform))

func CreateActor(world donburi.World, position math.Vec2) *donburi.Entry {
	entity := world.Create(ActorTag, Stats)
	entry := world.Entry(entity)
	Stats.SetValue(entry, *NewStatsData(StatsValue{
		StatMelee: 10.0,
		StatSpeed: 2.0,
	}))

	WithImage(entry, assets.RedSquareImage, position)
	WithCollision(entry)
	return entry
}

func EachActorAtPoint(world donburi.World, point image.Point, yield func(*donburi.Entry)) {
	for entry := range ActorQuery.Iter(world) {
		bounds, ok := Rect(entry)
		if !ok {
			continue
		}

		if point.In(bounds) {
			yield(entry)
		}
	}
}
