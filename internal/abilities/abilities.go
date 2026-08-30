package abilities

import (
	"rtwp_ebitengine/internal/systems"
	"uuid"
)

var AbilityRegistry = map[uuid.UUID]systems.Ability{
	Move.AbilityID: Move,
}
