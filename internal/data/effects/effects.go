package effects

import (
	"rtwp_ebitengine/internal/ecs"

	"github.com/google/uuid"
)

var EffectRegistry = map[uuid.UUID]ecs.Effect{
	// attack up
	AttackUp.EffectID: AttackUp,
}
