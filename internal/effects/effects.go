package effects

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/util"

	"github.com/google/uuid"
	"github.com/yohamta/donburi/ecs"
)

var EffectRegistry = map[uuid.UUID]entities.Effect{
	SpeedDown.EffectID: SpeedDown,
	SpeedUp.EffectID:   SpeedUp,
}

func ResolveModifiers(ecs *ecs.ECS, frame *util.Frame) {
	for modifier := range components.ModifierQuery.IterOrdered(ecs.World, components.Modifier) {
		instance := components.Modifier.Get(modifier)
		effect, ok := EffectRegistry[instance.EffectID]
		if ok {
			if effect.Active(ecs.World, modifier) {
				effect.Apply(ecs.World, frame, modifier)
			}
		}
	}
}
