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

	if player == nil || player.Action == nil || player.DragEnd != nil {
		screen.DrawImage(assets.CursorPointerImage, op)
		return
	}

	cursorPoint := util.NewVec2(x, y)
	data := player.Action.Data()
	if player.Action.Valid(ecs.World, cursorPoint) {
		if data.Cursor != nil {
			op.GeoM.Translate(data.CursorOffset.X, data.CursorOffset.Y)
			screen.DrawImage(data.Cursor, op)
		} else {
			screen.DrawImage(assets.CursorPointerImage, op)
		}
	} else {
		if data.CursorInvalid != nil {
			op.GeoM.Translate(data.CursorOffset.X, data.CursorOffset.Y)
			screen.DrawImage(data.CursorInvalid, op)
		} else {
		}
	}
}
