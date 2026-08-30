package effects

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

var systemModifiers = []components.Effect{
	SystemResolveStats,
}

func LoadSystemModifiers(ecs *ecs.ECS) {
	for _, sys := range systemModifiers {
		sys.Spawn(ecs, math.NewVec2(0, 0))
	}
}

func ResolveModifiers(ecs *ecs.ECS, frame *util.Frame) {
	for modifier := range components.ModifierQuery.IterOrdered(ecs.World, components.Modifier) {
		instance := components.Modifier.Get(modifier)
		effect := instance.Effect
		if effect.Active(ecs.World, modifier) {
			effect.Apply(ecs.World, frame, modifier)
		}
	}
}
