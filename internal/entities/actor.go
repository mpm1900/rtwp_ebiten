package entities

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

func CreateActor(ecs *ecs.ECS, position math.Vec2) donburi.Entity {
	entity := ecs.Create(ActorLayer, components.Actor, components.Stats)
	entry := ecs.World.Entry(entity)
	components.Stats.SetValue(entry, *components.NewStatsData(components.StatsValue{
		components.StatMelee: 10.0,
		components.StatSpeed: 2.0,
	}))

	components.WithImage(entry, assets.RedSquareImage, position)
	components.WithCollision(entry)
	return entity
}
