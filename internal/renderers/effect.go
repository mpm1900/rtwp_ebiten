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
		trans := transform.Transform.Get(entry)
		image := *components.Image.Get(entry)
		options := ebiten.DrawImageOptions{}

		center_scale := components.CenterScale(*trans)
		options.GeoM.Translate(center_scale.X, center_scale.Y)
		center := components.CenterTrans(*trans)
		center_point := view.Point(center)
		options.GeoM.Translate(center_point.X, center_point.Y)
		screen.DrawImage(image, &options)
	}
}
