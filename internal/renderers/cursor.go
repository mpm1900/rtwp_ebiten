package renderers

import (
	"rtwp_ebitengine/internal/abilities"
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
)

func RenderCursor(ecs *ecs.ECS, screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))

	player := components.GetPlayer(ecs.World)
	if player.Ability == nil || player.DragEnd != nil {
		screen.DrawImage(assets.CursorPointerImage, op)
		return
	}

	switch player.Ability.AbilityID {
	case abilities.Move.AbilityID:
		{
			if components.SelectedActorsQuery.Count(ecs.World) > 0 {
				screen.DrawImage(assets.CursorMoveImage, op)
				break
			} else {
				screen.DrawImage(assets.CursorPointerImage, op)
				break
			}
		}
	default:
		{
			screen.DrawImage(assets.CursorPointerImage, op)
		}
	}

}
