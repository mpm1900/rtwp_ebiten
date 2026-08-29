package systems

import (
	"rtwp_ebitengine/internal/events"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

func cursorPoint() math.Vec2 {
	x, y := ebiten.CursorPosition()
	return util.NewVec2(x, y)
}

func HandleSelection(ecs *ecs.ECS) {
	mousePoint := cursorPoint()

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
}
