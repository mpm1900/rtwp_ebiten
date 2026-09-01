package actions

import (
	"rtwp_ebitengine/internal/assets"
)

func LoadAbilities() {
	Move.Cursor = assets.CursorMoveImage
	Attack.Cursor = assets.CursorAttackImage
}
