package events

import (
	"image"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

var SelectInRect = events.NewEventType[image.Rectangle]()

func InitSelectInRect(world donburi.World) {
	SelectInRect.Subscribe(world, selectInRect)
}

func selectInRect(world donburi.World, rect image.Rectangle) {
	clearSelected(world, struct{}{})
	for entry := range components.ActorQuery.Iter(world) {
		actorRect, ok := components.Rect(entry)
		if !ok {
			continue
		}

		if rect.Overlaps(actorRect) {
			selectActor(world, entry.Entity())
		}
	}
}
