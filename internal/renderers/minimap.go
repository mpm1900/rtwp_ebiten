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
	camera := components.GetCamera(ecs.World)
	if camera.Camera == nil {
		return
	}

	worldRect := components.WorldRect()
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

	worldW := float64(worldRect.Dx())
	worldH := float64(worldRect.Dy())
	playerEntity := components.GetPlayerEntity(ecs.World)
	for entry := range components.ImageQuery.Iter(ecs.World) {
		position := components.Center(entry)
		mapX := (position.X - float64(worldRect.Min.X)) / worldW * components.MINIMAP_SIZE
		mapY := (position.Y - float64(worldRect.Min.Y)) / worldH * components.MINIMAP_SIZE
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

	scale := camera.Camera.Scale
	if scale <= 0 {
		scale = 1.0
	}

	halfWidth := float64(camera.Camera.Width) / 2.0 / scale
	halfHeight := float64(camera.Camera.Height) / 2.0 / scale

	viewLeft := camera.Camera.X - halfWidth
	viewTop := camera.Camera.Y - halfHeight
	viewRight := camera.Camera.X + halfWidth
	viewBottom := camera.Camera.Y + halfHeight

	viewLeft = max(viewLeft, float64(worldRect.Min.X))
	viewTop = max(viewTop, float64(worldRect.Min.Y))
	viewRight = min(viewRight, float64(worldRect.Max.X))
	viewBottom = min(viewBottom, float64(worldRect.Max.Y))

	viewportX := (viewLeft - float64(worldRect.Min.X)) / worldW * components.MINIMAP_SIZE
	viewportY := (viewTop - float64(worldRect.Min.Y)) / worldH * components.MINIMAP_SIZE
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
