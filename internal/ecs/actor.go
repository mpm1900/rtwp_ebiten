package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

var RedSquareImage *ebiten.Image
var BlueSquareImage *ebiten.Image
var GreenSquareImage *ebiten.Image

var ActorTag = donburi.NewTag("Actor")

func MakeActor(world donburi.World, position math.Vec2) *donburi.Entry {
	actor_entity := world.Create(ActorTag, Stats)
	actor_entry := world.Entry(actor_entity)
	Stats.SetValue(actor_entry, *NewStatsData(StatsValue{
		StatMelee: 10.0,
	}))

	WithImage(actor_entry, RedSquareImage, position)
	return actor_entry
}
