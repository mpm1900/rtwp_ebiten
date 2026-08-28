package components

import (
	"image"
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
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

func ActorBounds(actor *donburi.Entry) (image.Rectangle, bool) {
	if !actor.HasComponent(transform.Transform) || !actor.HasComponent(Image) {
		return image.Rectangle{}, false
	}

	trans := transform.Transform.Get(actor)
	sprite := *Image.Get(actor)
	if sprite == nil {
		return image.Rectangle{}, false
	}

	bounds := sprite.Bounds()
	return bounds.Add(util.ToPoint(trans.LocalPosition)), true
}
