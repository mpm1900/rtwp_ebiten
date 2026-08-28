package game

import (
	"image"
	"math"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.BeginDrag(mousePoint)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.UpdateDrag(mousePoint)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.EndDrag(mousePoint)
	}

	switch g.Action.Name {
	case ActionMove:
		{
			_, has_selection := components.Selected.First(g.World)
			if has_selection {
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
					g.MoveSelectedTo(mousePoint)
				}
			}
		}
	}
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
	push := ebiten.IsKeyPressed(ebiten.KeyShift)

	for selected := range components.Selected.Iter(g.World) {
		if push {
			components.PushMovement(selected, point)
			continue
		}
		components.WithMovement(selected, point)
	}
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

	rect := screenRect(*g.Selection.DragStart, *g.Selection.DragEnd)
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
