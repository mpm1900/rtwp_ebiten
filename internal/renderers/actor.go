package renderers

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var renderActorsQuery = donburi.NewQuery(filter.Contains(components.Actor, components.Image))

func RenderActors(ecs *ecs.ECS, screen *ebiten.Image) {
	for entry := range renderActorsQuery.Iter(ecs.World) {
		transform := transform.Transform.Get(entry)
		image := *components.Image.Get(entry)
		if entry.HasComponent(components.Selected) {
			image = assets.GreenSquareImage
		}

		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(transform.LocalPosition.X, transform.LocalPosition.Y)
		screen.DrawImage(image, options)
	}
}
