package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

var Actions = events.NewEventType[components.ActionEvent]()

func InitActions(world donburi.World) {
	Actions.Subscribe(world, handleActions)
}

func handleActions(world donburi.World, event components.ActionEvent) {
	event.Action.Handle(world, event.Source, event.Point)
}
