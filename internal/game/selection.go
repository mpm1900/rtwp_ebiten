package game

import (
	"image"
	"rtwp_ebitengine/internal/events"
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

func (g *Game) HandleSelection() {
	mousePoint := cursorPoint()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		events.ClearSelected.Publish(g.ECS.World, struct{}{})
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
		events.SelectAt.Publish(g.ECS.World, point)
	} else {
		dragRect, ok := g.DragRect()
		if ok {
			events.SelectInRect.Publish(g.ECS.World, dragRect)
		}
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
