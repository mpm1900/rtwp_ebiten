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

	defer AdvanceActionQueue(world, world.Entry(event.Source))
	return a.SetActionEvent(world, event)
}

func (a *ActorData) StartAction(world donburi.World, entry *donburi.Entry, event ActionEvent) bool {
	if a.CooldownForAction(event.Action) > 0 {
		return false
	}

	if entry.HasComponent(Delay) {
		if delay := *Delay.Get(entry); delay > 0 {
			return false
		}

		entry.RemoveComponent(Delay)
	} else if delay := event.Action.Delay; delay > 0 {
		WithDelay(entry, delay)
		return false
	}

	a.ActionStarted = true
	event.Action.Start(world, event)
	a.SetActionCooldown(event.Action)
	return true
}
func (a *ActorData) EndAction(world donburi.World, entry *donburi.Entry, event *ActionEvent) {
	if delay := event.Action.PostDelay; delay > 0 {
		WithDelay(entry, delay)
	}
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
