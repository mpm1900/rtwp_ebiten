package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
)

func HandleSelection(ecs *ecs.ECS) {
	mousePoint := util.CursorPoint()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		events.ClearSelected.Publish(ecs.World, struct{}{})
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		events.StartDrag.Publish(ecs.World, mousePoint)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		events.UpdateDrag.Publish(ecs.World, mousePoint)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		events.EndDrag.Publish(ecs.World, mousePoint)
	}

	player := components.GetPlayer(ecs.World)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		player.StartCameraDrag(mousePoint)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		delta, ok := player.UpdateCameraDrag(mousePoint)
		if ok {
			events.UpdateCamera.Publish(ecs.World, delta)
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		player.ClearCameraDrag()
	}
}
