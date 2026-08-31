package entities

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

func CreateActor(
	ecs *ecs.ECS,
	position math.Vec2,
	abilities []*components.Ability,
	player donburi.Entity,
) donburi.Entity {
	entity := ecs.Create(ActorLayer, components.Actor, components.Stats)
	entry := ecs.World.Entry(entity)
	components.Actor.SetValue(entry, components.ActorData{
		Abilities: abilities,
		Player:    player,
	})
	components.Stats.SetValue(entry, *components.NewStatsData(map[components.Stat]float64{
		components.StatHealth: 100.0,
		components.StatMelee:  10.0,
		components.StatSpeed:  2.0,
	}))
	components.WithDamage(entry, 0)
	components.WithImage(entry, assets.ActorImage, position)
	components.WithCollision(entry)
	return entity
}
