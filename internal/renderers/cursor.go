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
	base := &ebiten.DrawImageOptions{}
	base.GeoM.Translate(float64(x), float64(y))
	player := components.GetPlayer(ecs.World)

	if player.CameraDrag != nil {
		return
	}

	if player.Action == nil || player.DragEnd != nil {
		screen.DrawImage(assets.CursorPointerImage, base)
		return
	}

	mousePoint := util.CursorPoint()
	data := player.Action.Data()
	isValid := false

	if worldPoint, ok := components.MinimapWorldPoint(mousePoint); ok {
		isValid = player.Action.Valid(ecs.World, worldPoint)
	} else {
		isValid = player.Action.Valid(ecs.World, player.ScreenToWorld(mousePoint))
	}

	cursorImage := assets.CursorPointerImage
	if isValid {
		if data.Cursor != nil {
			cursorImage = data.Cursor
		}
	} else if data.CursorInvalid != nil {
		cursorImage = data.CursorInvalid
	}

	options := *base
	options.GeoM.Translate(data.CursorOffset.X, data.CursorOffset.Y)
	screen.DrawImage(cursorImage, &options)
}
