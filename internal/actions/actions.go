package actions

import (
	"rtwp_ebitengine/internal/assets"
)

func Load() {
	Move.Cursor = assets.CursorMoveImage
	Move.CursorInvalid = assets.CursorInvalidImage
	SingleAttack.Cursor = assets.CursorAttackImage
}
