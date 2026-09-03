package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

var Actions = events.NewEventType[components.ActionEvent]()
var ClearActions = events.NewEventType[struct{}]()

func InitActions(world donburi.World) {
	Actions.Subscribe(world, handleActions)
	ClearActions.Subscribe(world, handleClearActions)
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

		if !actor.ActionStarted {
			if actor.CooldownForAction(active_event.Action) > 0 {
				return
			}

			if entry.HasComponent(components.Delay) {
				delay := *components.Delay.Get(entry)
				if delay > 0 {
					return
				}

				entry.RemoveComponent(components.Delay)
			} else if delay := active_event.Action.Data().Delay; delay > 0 {
				components.WithDelay(entry, delay)
				return
			}

			event := *active_event
			actor.ActionStarted = true
			event.Action.Handle(world, event)
			actor.SetActionCooldown(event.Action)
			continue
		}

		if !active_event.Action.IsComplete(world, active_event.Source) {
			return
		}

		active_event.Action.Cancel(world, active_event.Source)
		actor.NextActionEvent()
	}
}

func handleClearActions(world donburi.World, _ struct{}) {
	for selected := range components.SelectedActorsQuery.Iter(world) {
		actor := components.Actor.Get(selected)
		action, ok := actor.PeekActionQueue()
		if ok {
			action.Action.Cancel(world, selected.Entity())
		}
		if selected.HasComponent(components.Delay) {
			selected.RemoveComponent(components.Delay)
		}
		actor.ActionStarted = false
		actor.ActionQueue.Clear()
	}
}
