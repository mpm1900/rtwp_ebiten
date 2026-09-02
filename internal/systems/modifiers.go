package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi/ecs"
)

func ResolveModifiers(frame *util.Frame) ecs.System {
	return func(ecs *ecs.ECS) {
		for modifier := range components.ModifierQuery.IterOrdered(ecs.World, components.Modifier) {
			instance := components.Modifier.Get(modifier)
			effect := instance.Effect
			if effect.Active(ecs.World, modifier) {
				effect.Apply(ecs.World, frame, modifier)
			}
		}
	}
}
