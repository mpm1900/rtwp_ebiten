package game

import (
	"image"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	dmath "github.com/yohamta/donburi/features/math"
)

const dragClickThreshold = 40

type SelectionState struct {
	DragStart *dmath.Vec2
	DragEnd   *dmath.Vec2
}

func NewSelection() SelectionState {
	return SelectionState{}
}

func (g *Game) ClearSelection() {
	for selected := range components.Selected.Iter(g.World) {
		selected.RemoveComponent(components.Selected)
	}
}

func (g *Game) HandleSelection() {
	mousePoint := cursorPoint()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.ClearSelection()
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.BeginDrag(mousePoint)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.UpdateDrag(mousePoint)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.EndDrag(mousePoint)
	}
}

func (g *Game) SelectAt(point dmath.Vec2) bool {
	mousePoint := util.ToPoint(point)
	g.ClearSelection()

	for entry := range components.ActorTag.Iter(g.World) {
		bounds, ok := components.ActorBounds(entry)
		if !ok {
			continue
		}

		if mousePoint.In(bounds) {
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
		actorRect, ok := components.ActorBounds(entry)
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

func (g *Game) BeginDrag(point dmath.Vec2) {
	g.Selection.DragStart = &point
	g.Selection.DragEnd = nil
}

func (g *Game) UpdateDrag(point dmath.Vec2) {
	if g.Selection.DragStart == nil {
		return
	}

	g.Selection.DragEnd = &point
}

func (g *Game) EndDrag(point dmath.Vec2) {
	if g.Selection.DragStart == nil {
		return
	}

	g.UpdateDrag(point)
	distance := g.Selection.DragStart.Distance(*g.Selection.DragEnd)
	if distance <= dragClickThreshold {
		g.SelectAt(*g.Selection.DragStart)
	} else {
		g.SelectInDragRect()
	}
	g.ClearDrag()
}

func (g *Game) ClearDrag() {
	g.Selection.DragStart = nil
	g.Selection.DragEnd = nil
}

func (g *Game) DragRect() (image.Rectangle, bool) {
	if g.Selection.DragStart == nil || g.Selection.DragEnd == nil {
		return image.Rectangle{}, false
	}

	rect := util.ToRect(*g.Selection.DragStart, *g.Selection.DragEnd)
	if rect.Dx() < 1 || rect.Dy() < 1 {
		return image.Rectangle{}, false
	}

	return rect, true
}
