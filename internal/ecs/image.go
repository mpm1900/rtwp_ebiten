package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var Image = donburi.NewComponentType[*ebiten.Image]()
var ImageQuery = donburi.NewQuery(filter.And(
	filter.Contains(Image),
	filter.Contains(Position),
))

func WithImage(entry *donburi.Entry, image *ebiten.Image, position Point) {
	entry.AddComponent(Image)
	Image.SetValue(entry, image)
	WithPosition(entry, position)
}

func RenderEntries(screen *ebiten.Image, world donburi.World) {
	for entry := range ImageQuery.Iter(world) {
		position := Position.Get(entry)
		image := *Image.Get(entry)
		if entry.HasComponent(Selected) {
			image = GreenSquareImage
		}

		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(position.X, position.Y)
		screen.DrawImage(image, options)
	}
}
