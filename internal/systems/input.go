package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

func handleMouseInput(ecs *ecs.ECS, point math.Vec2) {
	// zoom
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 && !components.IsOverMinimap(point, components.MinimapRect()) {
		events.ZoomCamera.Publish(ecs.World, events.ZoomCameraData{
			Delta:  wheelY,
			Cursor: point,
		})
	}

	// left
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if worldPoint, ok := components.MinimapWorldPoint(point); ok {
			events.LeftClickMinimap.Publish(ecs.World, worldPoint)
		} else {
			events.StartDrag.Publish(ecs.World, point)
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if worldPoint, ok := components.MinimapWorldPoint(point); ok {
			events.LeftClickMinimap.Publish(ecs.World, worldPoint)
		} else {
			events.UpdateDrag.Publish(ecs.World, point)
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		events.EndDrag.Publish(ecs.World, point)
	}

	player := components.GetPlayer(ecs.World)

	// right
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if worldPoint, ok := components.MinimapWorldPoint(point); ok {
			events.ActionClick.Publish(ecs.World, worldPoint)
		} else {
			events.ActionClick.Publish(ecs.World, player.ScreenToWorld(point))
		}
	}

	// middle
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		player.StartCameraDrag(point)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		delta, ok := player.UpdateCameraDrag(point)
		if ok {
			events.UpdateCamera.Publish(ecs.World, delta)
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		player.ClearCameraDrag()
	}
}

func HandleInput(ecs *ecs.ECS) {
	mousePoint := util.CursorPoint()
	handleMouseInput(ecs, mousePoint)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		events.ClearSelected.Publish(ecs.World, struct{}{})
	}
}
