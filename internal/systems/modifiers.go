package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi/ecs"
)

func ResolveModifiers(frame *util.Frame) ecs.System {
	return func(ecs *ecs.ECS) {
		for modifier := range components.ModifierQuery.IterOrdered(ecs.World, components.Modifier) {
			mod := components.Modifier.Get(modifier)
			if mod.Active(ecs.World, modifier) {
				mod.Apply(ecs.World, frame, modifier)
			}
		}
	}
}
