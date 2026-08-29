package entities

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlayer(esc *ecs.ECS) donburi.Entity {
	entity := esc.World.Create(components.Player)
	return entity
}
