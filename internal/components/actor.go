package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type ActorData struct {
	Player          donburi.Entity
	Actions         []*Action
	ActionQueue     util.Queue[ActionEvent]
	ActionCooldowns map[*Action]int
	ActionStarted   bool
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
		active_event.Action.Cancel(world, *active_event)
	}

	a.ActionStarted = false
	a.ActionQueue.Set(event)
	return true
}
func (a *ActorData) ReplaceActionEvent(event ActionEvent) bool {
	if !a.ActionQueue.Replace(event) {
		return a.PushActionEvent(event)
	}

	a.ActionStarted = false
	return true
}
func (a *ActorData) PushActionEvent(event ActionEvent) bool {
	should_start := a.ActionQueue.Push(event)
	if should_start {
		a.ActionStarted = false
	}

	return should_start
}
func (a *ActorData) PushNextActionEvent(event ActionEvent) bool {
	should_start := a.ActionQueue.InsertNext(event)
	if should_start {
		a.ActionStarted = false
	}

	return should_start
}
func (a *ActorData) QueueActionEvent(world donburi.World, event ActionEvent, push bool) bool {
	if push {
		return a.PushActionEvent(event)
	}

	return a.SetActionEvent(world, event)
}
func (a *ActorData) QueueAndAdvance(world donburi.World, entry *donburi.Entry, event ActionEvent, push bool) bool {
	if !a.QueueActionEvent(world, event, push) {
		return false
	}

	AdvanceActionQueue(world, entry)
	return true
}
func (a *ActorData) NextActionEvent() (*ActionEvent, bool) {
	next_event, ok := a.ActionQueue.Pop()
	a.ActionStarted = false
	return next_event, ok
}
func (a *ActorData) TickActionTimers() {
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
func (a *ActorData) CooldownForAction(action *Action) int {
	if action == nil {
		return 0
	}

	return a.ActionCooldowns[action]
}
func (a *ActorData) SetActionCooldown(action *Action) {
	if action == nil {
		return
	}

	cooldown := action.Cooldown
	if cooldown <= 0 {
		return
	}
	if a.ActionCooldowns == nil {
		a.ActionCooldowns = map[*Action]int{}
	}

	a.ActionCooldowns[action] = cooldown
}

var Actor = donburi.NewComponentType[ActorData]()
var ActorQuery = donburi.NewQuery(filter.Contains(Actor, transform.Transform))

func EachActorAtPoint(world donburi.World, point math.Vec2, yield func(*donburi.Entry)) {
	pt := image.Pt(int(point.X), int(point.Y))
	for entry := range ActorQuery.Iter(world) {
		bounds, ok := Rect(entry)
		if !ok {
			continue
		}

		if pt.In(bounds) {
			yield(entry)
		}
	}
}

func FirstActorAtPoint(world donburi.World, point math.Vec2) (*donburi.Entry, bool) {
	pt := image.Pt(int(point.X), int(point.Y))
	for entry := range ActorQuery.Iter(world) {
		bounds, ok := Rect(entry)
		if !ok {
			continue
		}

		if pt.In(bounds) {
			return entry, true
		}
	}

	return nil, false
}
