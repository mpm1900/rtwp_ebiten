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
		eventPoint := player.ScreenToWorld(point)
		if worldPoint, ok := components.MinimapWorldPoint(point); ok {
			eventPoint = worldPoint
		}
		keys := []ebiten.Key{}
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			keys = append(keys, ebiten.KeyShift)
		}
		if ebiten.IsKeyPressed(ebiten.KeyControl) {
			keys = append(keys, ebiten.KeyControl)
		}
		events.ActionClick.Publish(ecs.World, components.ActionEvent{
			Point: eventPoint,
			Keys:  keys,
		})
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

func handleActionInput(ecs *ecs.ECS) {
	player := components.GetPlayer(ecs.World)
	if player.SelectedAction == nil {
		return
	}

	actions := map[components.Action]int{}
	count := 0

	for selected := range components.SelectedActorsQuery.Iter(ecs.World) {
		count++
		actor := components.Actor.Get(selected)
		for _, action := range actor.Actions {
			actions[action]++
		}
	}

	for action := range actions {
		if inpututil.IsKeyJustPressed(action.Data().Key) && actions[action] == count {
			player.SelectedAction = action
		}
	}
}

func HandleInput(ecs *ecs.ECS) {
	mousePoint := util.CursorPoint()
	handleMouseInput(ecs, mousePoint)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		events.ClearSelected.Publish(ecs.World, struct{}{})
		events.ClearActions.Publish(ecs.World, struct{}{})
	}

	handleActionInput(ecs)
}
