package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type ActionEvent struct {
	Action           *Action
	Keys             []ebiten.Key
	Loop             bool
	Position         math.Vec2
	Source           donburi.Entity
	OnActionStart    func()
	OnActionComplete func()
}

type ActionStatus int

const (
	ActionRunning ActionStatus = iota
	ActionComplete
	ActionCanceled
)

type Action struct {
	ActionConfig
	Behavior ActionBehavior
}

type ActionConfig struct {
	Key           ebiten.Key
	Name          string
	Delay         int
	Cooldown      int
	Cursor        *ebiten.Image
	CursorInvalid *ebiten.Image
	CursorOffset  math.Vec2
	SingleOnly    bool
}

type ActionBehavior interface {
	// Publish turns input into queued ActionEvents for selected actors.
	Publish(*Action, donburi.World, ActionEvent)
	// Start applies one-time setup for a queued ActionEvent.
	Start(*Action, donburi.World, ActionEvent)
	// Update advances an active ActionEvent and reports its lifecycle state.
	Update(*Action, donburi.World, ActionEvent) ActionStatus
	// Cancel cleans up when the action is interrupted or popped.
	Cancel(*Action, donburi.World, ActionEvent)
	// Valid checks whether input at point can use this action.
	Valid(*Action, donburi.World, math.Vec2) bool
}

func (a *Action) Publish(world donburi.World, event ActionEvent) {
	if a == nil || a.Behavior == nil {
		return
	}

	a.Behavior.Publish(a, world, event)
}

func (a *Action) Start(world donburi.World, event ActionEvent) {
	if a == nil || a.Behavior == nil {
		return
	}

	if event.OnActionStart != nil {
		event.OnActionStart()
	}

	a.Behavior.Start(a, world, event)
}

func (a *Action) Update(world donburi.World, event ActionEvent) ActionStatus {
	if a == nil || a.Behavior == nil {
		return ActionComplete
	}

	return a.Behavior.Update(a, world, event)
}

func (a *Action) Cancel(world donburi.World, event ActionEvent) {
	if a == nil || a.Behavior == nil {
		return
	}

	if event.OnActionComplete != nil {
		event.OnActionComplete()
	}

	a.Behavior.Cancel(a, world, event)
}

func (a *Action) Valid(world donburi.World, point math.Vec2) bool {
	if a == nil || a.Behavior == nil {
		return false
	}

	return a.Behavior.Valid(a, world, point)
}
