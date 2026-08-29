package renderers

import (
	"image/color"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi/ecs"
)

func RenderDragRect(ecs *ecs.ECS, screen *ebiten.Image) {
	player := components.GetPlayer(ecs.World)
	if player.DragStart == nil || player.DragEnd == nil {
		return
	}
	rect := util.ToRect(*player.DragStart, *player.DragEnd)
	borderColor := color.RGBA{0, 0xff, 0, 0xff}
	vx := float32(rect.Min.X)
	vy := float32(rect.Min.Y)
	vheight := float32(rect.Dy())
	vwidth := float32(rect.Dx())
	vector.StrokeRect(screen, vx, vy, vwidth, vheight, 1, borderColor, false)
}
