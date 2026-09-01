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
	HandleActionQueue(world, event.Source)
}

func HandleActionQueue(world donburi.World, source donburi.Entity) {
	for {
		if !world.Valid(source) {
			return
		}

		entry := world.Entry(source)
		if !entry.HasComponent(components.Actor) {
			return
		}

		actor := components.Actor.Get(entry)
		active_event, ok := actor.PeekActionQueue()
		if !ok {
			return
		}
		if active_event.Action == nil {
			actor.NextActionEvent()
			continue
		}

		if !active_event.Started {
			active_event.Started = true
			active_event.Action.Handle(world, active_event.Source, active_event.Point)
			continue
		}

		if !active_event.Action.IsComplete(world, active_event.Source) {
			return
		}

		actor.NextActionEvent()
	}
}
