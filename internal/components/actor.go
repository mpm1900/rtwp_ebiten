package components

import (
	"image"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type ActorData struct {
	Player      donburi.Entity
	Actions     []Action
	ActionQueue []ActionEvent
}

func (a *ActorData) PeekActionQueue() (*ActionEvent, bool) {
	if len(a.ActionQueue) == 0 {
		return nil, false
	}

	return &a.ActionQueue[0], true
}
func (a *ActorData) HasAction() bool {
	return len(a.ActionQueue) > 0
}
func (a *ActorData) SetActionEvent(world donburi.World, event ActionEvent) bool {
	if active_event, ok := a.PeekActionQueue(); ok && active_event.Started {
		active_event.Action.Cancel(world, active_event.Source)
	}

	event.Started = false
	clear(a.ActionQueue)
	a.ActionQueue = append(a.ActionQueue[:0], event)
	return true
}
func (a *ActorData) PushActionEvent(event ActionEvent) bool {
	event.Started = false
	a.ActionQueue = append(a.ActionQueue, event)
	return len(a.ActionQueue) == 1
}
func (a *ActorData) NextActionEvent() (*ActionEvent, bool) {
	if len(a.ActionQueue) == 0 {
		return nil, false
	}

	copy(a.ActionQueue, a.ActionQueue[1:])
	last_index := len(a.ActionQueue) - 1
	a.ActionQueue[last_index] = ActionEvent{}
	a.ActionQueue = a.ActionQueue[:last_index]
	return a.PeekActionQueue()
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
