package actions

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type NullAction struct {
	components.ActionData
}

func (a NullAction) Data() components.ActionData {
	return a.ActionData
}
func (a NullAction) Handle(world donburi.World, point math.Vec2) {
}
func (a NullAction) Valid(world donburi.World, point math.Vec2) bool {
	return false
}

var Null = NullAction{}

func LoadAbilities() {
	Null.CursorInvalid = assets.CursorPointerImage
	Move.Cursor = assets.CursorMoveImage
	Attack.Cursor = assets.CursorAttackImage
}
