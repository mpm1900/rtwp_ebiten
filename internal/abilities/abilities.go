package abilities

import (
	"rtwp_ebitengine/internal/components"

	"github.com/google/uuid"
)

var AbilityRegistry = map[uuid.UUID]components.Ability{
	Move.AbilityID: Move,
}
