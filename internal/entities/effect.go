package entities

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateEffect(ecs *ecs.ECS, effect components.Effect, mod components.ModifierConfig) donburi.Entity {
	return CreateModifier(ecs, components.NewModifierInstance(effect, mod), EffectLayer)
}

func CreateModifier(ecs *ecs.ECS, mod components.ModifierData, layer ecs.LayerID) donburi.Entity {
	entity := ecs.Create(layer, components.Modifier)
	entry := ecs.World.Entry(entity)
	components.Modifier.SetValue(entry, mod)
	return entity
}
