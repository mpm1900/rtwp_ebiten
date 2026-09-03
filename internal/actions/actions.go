package actions

import (
	"rtwp_ebitengine/internal/assets"
)

func LoadAbilities() {
	Move.Cursor = assets.CursorMoveImage
	Move.CursorInvalid = assets.CursorInvalidImage
	Attack.Cursor = assets.CursorAttackImage
}
