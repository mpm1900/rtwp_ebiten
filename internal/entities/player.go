package entities

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlayer(esc *ecs.ECS, startingAbility *components.Ability) donburi.Entity {
	entity := esc.World.Create(components.Player)
	entry := esc.World.Entry(entity)
	components.Player.SetValue(entry, components.NewPlayerData(startingAbility))
	return entity
}
