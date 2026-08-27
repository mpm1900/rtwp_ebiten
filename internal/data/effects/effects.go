package effects

import (
	"rtwp_ebitengine/internal/components"

	"github.com/google/uuid"
)

var EffectRegistry = map[uuid.UUID]components.Effect{
	// attack up
	SpeedUp.EffectID: SpeedUp,
}
