package entities

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

type Effect interface {
	Modifier() components.ModifierData
	Active(world donburi.World, modifier *donburi.Entry) bool
	Apply(world donburi.World, frame *util.Frame, modifier *donburi.Entry)
	Spawn(ecs *ecs.ECS, position math.Vec2) donburi.Entity
}

func CreateEffect(ecs *ecs.ECS, effect Effect) donburi.Entity {
	return CreateModifier(ecs, effect.Modifier(), EffectLayer)
}

func CreateModifier(ecs *ecs.ECS, mod components.ModifierData, layer ecs.LayerID) donburi.Entity {
	entity := ecs.Create(layer, components.Modifier)
	entry := ecs.World.Entry(entity)
	components.Modifier.SetValue(entry, mod)
	return entity
}
