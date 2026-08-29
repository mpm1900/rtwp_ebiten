package components

import (
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
	filter.Not(filter.Contains(Delay)),
))

func WithImage(entry *donburi.Entry, image *ebiten.Image, position math.Vec2) {
	entry.AddComponent(Image)
	Image.SetValue(entry, image)

	scale := math.NewVec2(1, 1)
	if image != nil {
		bounds := image.Bounds()
		scale = math.NewVec2(float64(bounds.Dx()), float64(bounds.Dy()))
	}

	WithTransform(entry, transform.TransformData{
		LocalPosition: position,
		LocalScale:    scale,
	})
}
