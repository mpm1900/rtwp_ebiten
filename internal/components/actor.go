package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type QueuedActionEvent struct {
	ActionEvent
	DelayStarted bool
	Started      bool
}
type ActorData struct {
	Player         donburi.Entity
	Actions        []Action
	ActionQueue    util.Queue[QueuedActionEvent]
	ActionDelay    int
	ActionCooldown int
}

func NewQueuedActionEvent(e ActionEvent) QueuedActionEvent {
	return QueuedActionEvent{
		ActionEvent: e,
		Started:     false,
	}
}

func (a *ActorData) PeekActionQueue() (*QueuedActionEvent, bool) {
	return a.ActionQueue.Peek()
}
func (a *ActorData) HasAction() bool {
	return a.ActionQueue.Len() > 0
}
func (a *ActorData) ActionQueueLen() int {
	return a.ActionQueue.Len()
}
func (a *ActorData) SetActionEvent(world donburi.World, event QueuedActionEvent) bool {
	if active_event, ok := a.PeekActionQueue(); ok && active_event.Started {
		active_event.Action.Cancel(world, active_event.Source)
	}

	event.DelayStarted = false
	event.Started = false
	a.ActionDelay = 0
	a.ActionCooldown = 0
	a.ActionQueue.Set(event)
	return true
}
func (a *ActorData) PushActionEvent(event QueuedActionEvent) bool {
	event.DelayStarted = false
	event.Started = false
	return a.ActionQueue.Push(event)
}
func (a *ActorData) QueueActionEvent(world donburi.World, event ActionEvent, push bool) bool {
	if push {
		return a.PushActionEvent(NewQueuedActionEvent(event))
	}

	return a.SetActionEvent(world, NewQueuedActionEvent(event))
}
func (a *ActorData) NextActionEvent() (*QueuedActionEvent, bool) {
	return a.ActionQueue.Pop()
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
