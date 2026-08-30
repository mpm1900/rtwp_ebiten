package components

import (
	"image"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type ActorData struct {
	Abilities []*Ability
}

var Actor = donburi.NewComponentType[ActorData]()
var ActorQuery = donburi.NewQuery(filter.Contains(Actor, transform.Transform))

func EachActorAtPoint(world donburi.World, point image.Point, yield func(*donburi.Entry)) {
	for entry := range ActorQuery.Iter(world) {
		bounds, ok := Rect(entry)
		if !ok {
			continue
		}

		if point.In(bounds) {
			yield(entry)
		}
	}
}

func FirstActorAtPoint(world donburi.World, point image.Point) (*donburi.Entry, bool) {
	for entry := range ActorQuery.Iter(world) {
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
