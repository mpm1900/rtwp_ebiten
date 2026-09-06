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
	cursor := components.GetCursor(ecs.World)
	camera := components.GetCamera(ecs.World)
	player := components.GetPlayer(ecs.World)

	if camera.CameraDrag != nil {
		return
	}

	if player.SelectedAction == nil || cursor.DragEnd != nil {
		screen.DrawImage(assets.CursorPointerImage, base)
		return
	}

	mousePoint := util.CursorPoint()
	action := player.SelectedAction
	isValid := false

	if worldPoint, ok := components.MinimapWorldPoint(mousePoint); ok {
		isValid = action.Valid(ecs.World, worldPoint)
	} else {
		isValid = action.Valid(ecs.World, camera.ScreenToWorld(mousePoint))
	}

	var cursorImage *ebiten.Image
	if isValid {
		if action.Cursor != nil {
			cursorImage = action.Cursor
		}
	} else if action.CursorInvalid != nil {
		cursorImage = action.CursorInvalid
	}

	if cursorImage != nil {
		options := *base
		options.GeoM.Translate(action.CursorOffset.X, action.CursorOffset.Y)
		screen.DrawImage(cursorImage, &options)
	}
}
