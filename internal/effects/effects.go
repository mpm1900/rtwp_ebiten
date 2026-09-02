package effects

import (
	"rtwp_ebitengine/internal/components"

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
