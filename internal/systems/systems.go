package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/google/uuid"
	"github.com/yohamta/donburi/ecs"
)

func Load(ecs *ecs.ECS, frame *util.Frame, abilityRegistry map[uuid.UUID]*components.Ability) {
	ecs.AddSystem(DecrementDelays)
	ecs.AddSystem(RemoveCompletedDelays)
	ecs.AddSystem(DecrementDurations)
	ecs.AddSystem(RemoveCompleted)
	ecs.AddSystem(ResolveModifiers(frame))

	ecs.AddSystem(HandleSelection)
	ecs.AddSystem(HandleAbilities(abilityRegistry))
	ecs.AddSystem(MoveEntities)
}
