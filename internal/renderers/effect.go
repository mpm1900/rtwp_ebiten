package renderers

import (
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var renderEffectsQuery = donburi.NewQuery(filter.And(
	filter.Contains(components.Modifier, components.Image),
	filter.Not(filter.Contains(components.Delay)),
))

func RenderEffect(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)

	for entry := range renderEffectsQuery.Iter(ecs.World) {
		transform := transform.Transform.Get(entry)
		image := *components.Image.Get(entry)
		options := ebiten.DrawImageOptions{}
		viewpoint := view.Point(transform.LocalPosition)
		options.GeoM.Translate(viewpoint.X, viewpoint.Y)
		screen.DrawImage(image, &options)
	}
}
