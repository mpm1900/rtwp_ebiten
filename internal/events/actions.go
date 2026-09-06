package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

var ClearActions = events.NewEventType[struct{}]()

func InitActions(world donburi.World) {
	ClearActions.Subscribe(world, handleClearActions)
}

func handleClearActions(world donburi.World, _ struct{}) {
	for selected := range components.SelectedActorsQuery.Iter(world) {
		actor := components.Actor.Get(selected)
		action, ok := actor.PeekActionQueue()
		if ok {
			action.Action.Cancel(world, *action)
		}
		if selected.HasComponent(components.Delay) {
			selected.RemoveComponent(components.Delay)
		}
		actor.ActionStarted = false
		actor.ActionQueue.Clear()
	}
}
