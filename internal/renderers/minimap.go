package renderers

import (
	"image/color"
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi/ecs"
)

func RenderMinimap(ecs *ecs.ECS, screen *ebiten.Image) {
	player := components.GetPlayer(ecs.World)
	if player == nil || player.Camera == nil {
		return
	}

	worldMinX, worldMinY, worldMaxX, worldMaxY := components.WorldRect()
	mapRect := components.MinimapRect()

	vector.StrokeRect(
		screen,
		float32(mapRect.Min.X-components.MINIMAP_BORDER),
		float32(mapRect.Min.Y-components.MINIMAP_BORDER),
		float32(components.MINIMAP_SIZE+components.MINIMAP_BORDER*2),
		float32(components.MINIMAP_SIZE+components.MINIMAP_BORDER*2),
		1,
		color.Black,
		false,
	)
	vector.FillRect(
		screen,
		float32(mapRect.Min.X),
		float32(mapRect.Min.Y),
		components.MINIMAP_SIZE,
		components.MINIMAP_SIZE,
		assets.ColorMinimap,
		false,
	)

	worldW := worldMaxX - worldMinX
	worldH := worldMaxY - worldMinY
	viewLeft := player.Camera.X - float64(player.Camera.Width)/2.0
	viewTop := player.Camera.Y - float64(player.Camera.Height)/2.0
	viewRight := player.Camera.X + float64(player.Camera.Width)/2.0
	viewBottom := player.Camera.Y + float64(player.Camera.Height)/2.0

	viewLeft = max(viewLeft, worldMinX)
	viewTop = max(viewTop, worldMinY)
	viewRight = min(viewRight, worldMaxX)
	viewBottom = min(viewBottom, worldMaxY)

	viewportX := (viewLeft - worldMinX) / worldW * components.MINIMAP_SIZE
	viewportY := (viewTop - worldMinY) / worldH * components.MINIMAP_SIZE
	viewportW := (viewRight - viewLeft) / worldW * components.MINIMAP_SIZE
	viewportH := (viewBottom - viewTop) / worldH * components.MINIMAP_SIZE

	vector.StrokeRect(
		screen,
		float32(mapRect.Min.X)+float32(viewportX),
		float32(mapRect.Min.Y)+float32(viewportY),
		float32(viewportW),
		float32(viewportH),
		2.0,
		assets.ColorSelected,
		false,
	)
}
