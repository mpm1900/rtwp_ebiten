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
	playerEntity := components.GetPlayerEntity(ecs.World)
	for entry := range components.ImageQuery.Iter(ecs.World) {
		position := components.Center(entry)
		mapX := (position.X - worldMinX) / worldW * components.MINIMAP_SIZE
		mapY := (position.Y - worldMinY) / worldH * components.MINIMAP_SIZE
		if mapX < 0 || mapX >= components.MINIMAP_SIZE || mapY < 0 || mapY >= components.MINIMAP_SIZE {
			continue
		}

		var markerColor color.Color = color.White
		if entry.HasComponent(components.Actor) {
			markerColor = assets.ColorEnemy
			actor := components.Actor.Get(entry)
			if actor.Player == playerEntity {
				markerColor = assets.ColorPlayer
			}
		}

		screen.Set(
			mapRect.Min.X+int(mapX),
			mapRect.Min.Y+int(mapY),
			markerColor,
		)
	}

	scale := player.Camera.Scale
	if scale <= 0 {
		scale = 1.0
	}

	halfWidth := float64(player.Camera.Width) / 2.0 / scale
	halfHeight := float64(player.Camera.Height) / 2.0 / scale

	viewLeft := player.Camera.X - halfWidth
	viewTop := player.Camera.Y - halfHeight
	viewRight := player.Camera.X + halfWidth
	viewBottom := player.Camera.Y + halfHeight

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
