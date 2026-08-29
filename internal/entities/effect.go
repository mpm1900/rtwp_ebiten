package entities

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Effect interface {
	Modifier() components.ModifierData
	Active(world donburi.World, modifier *donburi.Entry) bool
	Apply(world donburi.World, frame *util.Frame, modifier *donburi.Entry)
}

func CreateEffect(ecs *ecs.ECS, effect Effect, layer ecs.LayerID) *donburi.Entry {
	return CreateModifier(ecs, effect.Modifier(), layer)
}

func CreateModifier(ecs *ecs.ECS, mod components.ModifierData, layer ecs.LayerID) *donburi.Entry {
	entity := ecs.Create(layer, components.Modifier)
	entry := ecs.World.Entry(entity)
	components.Modifier.SetValue(entry, mod)
	return entry
}
