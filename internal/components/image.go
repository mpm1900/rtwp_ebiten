package components

import (
	"rtwp_ebitengine/internal/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var Image = donburi.NewComponentType[*ebiten.Image]()
var ImageQuery = donburi.NewQuery(filter.And(
	filter.Contains(Image),
	filter.Contains(transform.Transform),
))

func WithImage(entry *donburi.Entry, image *ebiten.Image, position math.Vec2) {
	entry.AddComponent(Image)
	Image.SetValue(entry, image)

	entry.AddComponent(transform.Transform)
	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: position,
	})
}

func RenderEntries(screen *ebiten.Image, world donburi.World) {
	for entry := range ImageQuery.Iter(world) {
		transform := transform.Transform.Get(entry)
		image := *Image.Get(entry)
		if entry.HasComponent(Selected) {
			image = assets.GreenSquareImage
		}

		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(transform.LocalPosition.X, transform.LocalPosition.Y)
		screen.DrawImage(image, options)
	}
}
