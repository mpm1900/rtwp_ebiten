package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type ActorData struct {
	Player          donburi.Entity
	Actions         []Action
	ActionQueue     util.Queue[ActionEvent]
	ActionDelay     int
	ActionCooldowns map[ebiten.Key]int
	ActionStarted   bool
	DelayStarted    bool
}

func (a *ActorData) PeekActionQueue() (*ActionEvent, bool) {
	return a.ActionQueue.Peek()
}
func (a *ActorData) HasAction() bool {
	return a.ActionQueue.Len() > 0
}
func (a *ActorData) ActionQueueLen() int {
	return a.ActionQueue.Len()
}
func (a *ActorData) SetActionEvent(world donburi.World, event ActionEvent) bool {
	if active_event, ok := a.PeekActionQueue(); ok && a.ActionStarted {
		active_event.Action.Cancel(world, active_event.Source)
	}

	a.DelayStarted = false
	a.ActionStarted = false
	a.ActionDelay = 0
	a.ActionQueue.Set(event)
	return true
}
func (a *ActorData) PushActionEvent(event ActionEvent) bool {
	should_start := a.ActionQueue.Push(event)
	if should_start {
		a.DelayStarted = false
		a.ActionStarted = false
		a.ActionDelay = 0
	}

	return should_start
}
func (a *ActorData) QueueActionEvent(world donburi.World, event ActionEvent, push bool) bool {
	if push {
		return a.PushActionEvent(event)
	}

	return a.SetActionEvent(world, event)
}
func (a *ActorData) NextActionEvent() (*ActionEvent, bool) {
	next_event, ok := a.ActionQueue.Pop()
	a.DelayStarted = false
	a.ActionStarted = false
	a.ActionDelay = 0
	return next_event, ok
}
func (a *ActorData) TickActionTimers() {
	if a.ActionDelay > 0 {
		a.ActionDelay--
	}

	for action_key, cooldown := range a.ActionCooldowns {
		cooldown--
		if cooldown <= 0 {
			delete(a.ActionCooldowns, action_key)
			continue
		}

		a.ActionCooldowns[action_key] = cooldown
	}
}
func (a *ActorData) HasActionCooldowns() bool {
	return len(a.ActionCooldowns) > 0
}
func (a *ActorData) CooldownForAction(action Action) int {
	if action == nil {
		return 0
	}

	return a.ActionCooldowns[action.Data().Key]
}
func (a *ActorData) SetActionCooldown(action Action) {
	if action == nil {
		return
	}

	cooldown := action.Data().Cooldown
	if cooldown <= 0 {
		return
	}
	if a.ActionCooldowns == nil {
		a.ActionCooldowns = map[ebiten.Key]int{}
	}

	a.ActionCooldowns[action.Data().Key] = cooldown
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
