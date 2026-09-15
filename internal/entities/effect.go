package entities

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateModifier(ecs *ecs.ECS, mod components.ModifierData) donburi.Entity {
	entity := ecs.Create(EffectLayer, components.Modifier)
	entry := ecs.World.Entry(entity)
	components.Modifier.SetValue(entry, mod)
	return entity
}
