package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

var ClearSelected = events.NewEventType[struct{}]()

func InitClearSelected(world donburi.World) {
	ClearSelected.Subscribe(world, clearSelected)
}

func clearSelected(world donburi.World, _ struct{}) {
	for selected := range components.Selected.Iter(world) {
		selected.RemoveComponent(components.Selected)
	}
}
