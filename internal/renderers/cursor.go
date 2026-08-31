package renderers

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
)

func RenderCursor(ecs *ecs.ECS, screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	player := components.GetPlayer(ecs.World)

	if player == nil || player.Ability == nil || player.DragEnd != nil {
		screen.DrawImage(assets.CursorPointerImage, op)
		return
	}

	cursorPoint := util.NewVec2(x, y)
	if player.Ability.Valid != nil && player.Ability.Valid(ecs, cursorPoint) {
		if player.Ability.Cursor != nil {
			op.GeoM.Translate(player.Ability.CursorOffset.X, player.Ability.CursorOffset.Y)
			screen.DrawImage(player.Ability.Cursor, op)
		} else {
			screen.DrawImage(assets.CursorPointerImage, op)
		}
	} else {
		if player.Ability.CursorInvalid != nil {
			op.GeoM.Translate(player.Ability.CursorOffset.X, player.Ability.CursorOffset.Y)
			screen.DrawImage(player.Ability.CursorInvalid, op)
		} else {
		}
	}
}
