package game

import (
	"image"
	"math"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

const dragClickThreshold = 40

type State struct {
	DragStart *dmath.Vec2
	DragEnd   *dmath.Vec2
}

func NewState() State {
	return State{}
}

func (g *Game) ClearSelection() {
	for selected := range components.Selected.Iter(g.World) {
		selected.RemoveComponent(components.Selected)
	}
}

func (g *Game) HandleClick(point dmath.Vec2) {
	if g.SelectAt(point) {
		return
	}

	g.MoveSelectedTo(point)
}

func (g *Game) SelectAt(point dmath.Vec2) bool {
	mousePoint := screenPoint(point)

	for entry := range components.ActorTag.Iter(g.World) {
		bounds, ok := actorBounds(entry)
		if !ok {
			continue
		}

		if mousePoint.In(bounds) {
			g.ClearSelection()
			entry.AddComponent(components.Selected)
			return true
		}
	}

	return false
}

func (g *Game) SelectInDragRect() bool {
	valid := false
	dragRect, ok := g.DragRect()
	if !ok {
		return false
	}

	g.ClearSelection()
	for entry := range components.ActorTag.Iter(g.World) {
		actorRect, ok := actorBounds(entry)
		if !ok {
			continue
		}

		if dragRect.Overlaps(actorRect) {
			entry.AddComponent(components.Selected)
			valid = true
		}
	}

	return valid
}

func (g *Game) MoveSelectedTo(point dmath.Vec2) {
	for selected := range components.Selected.Iter(g.World) {
		components.WithMovement(selected, point)
	}
}

func (g *Game) BeginDrag(point dmath.Vec2) {
	g.State.DragStart = &point
	g.State.DragEnd = nil
}

func (g *Game) UpdateDrag(point dmath.Vec2) {
	if g.State.DragStart == nil {
		return
	}

	g.State.DragEnd = &point
}

func (g *Game) EndDrag(point dmath.Vec2) {
	if g.State.DragStart == nil {
		return
	}

	g.UpdateDrag(point)
	distance := g.State.DragStart.Distance(*g.State.DragEnd)
	if distance <= dragClickThreshold {
		g.HandleClick(*g.State.DragStart)
	} else {
		g.SelectInDragRect()
	}

	g.ClearDrag()
}

func (g *Game) ClearDrag() {
	g.State.DragStart = nil
	g.State.DragEnd = nil
}

func (g *Game) DragRect() (image.Rectangle, bool) {
	if g.State.DragStart == nil || g.State.DragEnd == nil {
		return image.Rectangle{}, false
	}

	rect := screenRect(*g.State.DragStart, *g.State.DragEnd)
	if rect.Dx() < 1 || rect.Dy() < 1 {
		return image.Rectangle{}, false
	}

	return rect, true
}

func actorBounds(entry *donburi.Entry) (image.Rectangle, bool) {
	if !entry.HasComponent(transform.Transform) || !entry.HasComponent(components.Image) {
		return image.Rectangle{}, false
	}

	trans := transform.Transform.Get(entry)
	sprite := *components.Image.Get(entry)
	if sprite == nil {
		return image.Rectangle{}, false
	}

	bounds := sprite.Bounds()
	min := screenPoint(trans.LocalPosition)
	max := image.Pt(
		int(math.Ceil(trans.LocalPosition.X+float64(bounds.Dx()))),
		int(math.Ceil(trans.LocalPosition.Y+float64(bounds.Dy()))),
	)
	return image.Rectangle{
		Min: min,
		Max: max,
	}, true
}

func screenPoint(point dmath.Vec2) image.Point {
	return image.Pt(
		int(math.Floor(point.X)),
		int(math.Floor(point.Y)),
	)
}

func screenRect(start, end dmath.Vec2) image.Rectangle {
	min := image.Pt(
		int(math.Floor(math.Min(start.X, end.X))),
		int(math.Floor(math.Min(start.Y, end.Y))),
	)
	max := image.Pt(
		int(math.Ceil(math.Max(start.X, end.X))),
		int(math.Ceil(math.Max(start.Y, end.Y))),
	)
	return image.Rectangle{
		Min: min,
		Max: max,
	}
}
