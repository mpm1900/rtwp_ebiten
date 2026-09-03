package components

import (
	"image"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type InteractableData struct {
	TargetOffset math.Vec2
}

func (i InteractableData) Point(p math.Vec2) math.Vec2 {
	return p.Add(i.TargetOffset)
}

var Interactable = donburi.NewComponentType[InteractableData]()
var InteractableQuery = donburi.NewQuery(
	filter.And(
		filter.Contains(Interactable),
		filter.Not(filter.Contains(Delay)),
	))

func WitherInteractable(entry *donburi.Entry, data InteractableData) {
	entry.AddComponent(Interactable)
	Interactable.SetValue(entry, data)
	if !entry.HasComponent(transform.Transform) {
		WithTransform(entry, transform.TransformData{})
	}
}

func EachInteractableAtPoint(world donburi.World, point image.Point, yield func(*donburi.Entry)) {
	for entry := range Interactable.Iter(world) {
		bounds, ok := Rect(entry)
		if !ok {
			continue
		}

		if point.In(bounds) {
			yield(entry)
		}
	}
}

func FirstInteractableAtPoint(world donburi.World, point image.Point) (*donburi.Entry, bool) {
	for entry := range Interactable.Iter(world) {
		bounds, ok := Rect(entry)
		if !ok {
			continue
		}

		if point.In(bounds) {
			return entry, true
		}
	}

	return nil, false
}
