package components

import (
	"rtwp_ebitengine/internal/assets"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

var ActorTag = donburi.NewTag("Actor")

func CreateActor(world donburi.World, position math.Vec2) *donburi.Entry {
	entity := world.Create(ActorTag, Stats)
	entry := world.Entry(entity)
	Stats.SetValue(entry, *NewStatsData(StatsValue{
		StatMelee: 10.0,
		StatSpeed: 2.0,
	}))

	WithImage(entry, assets.RedSquareImage, position)
	WithCollision(entry, math.NewVec2(24, 24))
	return entry
}
