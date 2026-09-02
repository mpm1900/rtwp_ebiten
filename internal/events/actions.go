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
			if !actor.DelayStarted {
				actor.DelayStarted = true
				actor.ActionDelay = active_event.Action.Data().Delay
			}
			if actor.ActionDelay > 0 {
				return
			}

			actor.ActionStarted = true
			active_event.Action.Handle(world, *active_event)
			actor.SetActionCooldown(active_event.Action)
			continue
		}

		if !active_event.Action.IsComplete(world, active_event.Source) {
			return
		}

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
		actor.ActionDelay = 0
		actor.ActionStarted = false
		actor.DelayStarted = false
		actor.ActionQueue.Clear()
	}
}
