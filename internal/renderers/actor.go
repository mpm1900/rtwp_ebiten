package renderers

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var renderActorsQuery = donburi.NewQuery(filter.Contains(components.Actor, components.Image))

func RenderActors(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)

	for entry := range renderActorsQuery.Iter(ecs.World) {
		trans := transform.Transform.Get(entry)
		image := *components.Image.Get(entry)

		if entry.HasComponent(components.Selected) {
			image = assets.GreenSquareImage
		}

		options := ebiten.DrawImageOptions{}

		centerScale := components.CenterScale(*trans)
		options.GeoM.Translate(centerScale.X, centerScale.Y)

		options.GeoM.Rotate(dmath.ToRadians(trans.LocalRotation))

		center := components.CenterTrans(*trans)
		centerPoint := view.Point(center)
		options.GeoM.Translate(centerPoint.X, centerPoint.Y)

		screen.DrawImage(image, &options)
	}
}
