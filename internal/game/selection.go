package game

import (
	"fmt"
	"math"
	"rtwp_ebitengine/internal/ecs"

	"github.com/yohamta/donburi"
)

const DRAG_CLICK_THRESHOLD = 40

type State struct {
	DragStart *ecs.Point
	DragEnd   *ecs.Point
}

func NewState() State {
	return State{}
}

func (g *Game) ClearSelection() {
	for selected := range ecs.Selected.Iter(g.World) {
		selected.RemoveComponent(ecs.Selected)
	}
}

func (g *Game) HandleClick(point ecs.Point) {
	if g.SelectAt(point) {
		return
	}

	g.MoveSelectedTo(point)
}

func (g *Game) SelectAt(point ecs.Point) bool {
	for entry := range ecs.ActorTag.Iter(g.World) {
		x, y, width, height, ok := actorBounds(entry)
		if !ok {
			continue
		}

		if pointInRect(point, x, y, width, height) {
			g.ClearSelection()
			entry.AddComponent(ecs.Selected)
			return true
		}
	}

	return false
}

func (g *Game) SelectInDragRect() bool {
	valid := false
	x, y, width, height, ok := g.DragRect()
	if !ok {
		return false
	}

	g.ClearSelection()
	for entry := range ecs.ActorTag.Iter(g.World) {
		actorX, actorY, actorWidth, actorHeight, ok := actorBounds(entry)
		if !ok {
			continue
		}

		if rectsOverlap(x, y, width, height, actorX, actorY, actorWidth, actorHeight) {
			entry.AddComponent(ecs.Selected)
			valid = true
		}
	}

	return valid
}

func (g *Game) MoveSelectedTo(point ecs.Point) {
	for selected := range ecs.Selected.Iter(g.World) {
		ecs.WithMovement(selected, point)
	}
}

func (g *Game) BeginDrag(point ecs.Point) {
	g.State.DragStart = &point
	g.State.DragEnd = nil
}

func (g *Game) UpdateDrag(point ecs.Point) {
	if g.State.DragStart == nil {
		return
	}

	g.State.DragEnd = &point
}

func (g *Game) EndDrag(point ecs.Point) {
	if g.State.DragStart == nil {
		return
	}

	g.UpdateDrag(point)
	distance := g.DragDistance()
	fmt.Println(distance)
	if distance <= DRAG_CLICK_THRESHOLD {
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

func (g *Game) DragDistance() float64 {
	if g.State.DragStart == nil || g.State.DragEnd == nil {
		return 0
	}

	dx := g.State.DragEnd.X - g.State.DragStart.X
	dy := g.State.DragEnd.Y - g.State.DragStart.Y
	return math.Hypot(dx, dy)
}

func (g *Game) DragRect() (x, y, width, height float64, ok bool) {
	if g.State.DragStart == nil || g.State.DragEnd == nil {
		return 0, 0, 0, 0, false
	}

	start := g.State.DragStart
	end := g.State.DragEnd
	x = math.Min(start.X, end.X)
	y = math.Min(start.Y, end.Y)
	width = math.Abs(end.X - start.X)
	height = math.Abs(end.Y - start.Y)
	if width < 1 || height < 1 {
		return 0, 0, 0, 0, false
	}

	return x, y, width, height, true
}

func actorBounds(entry *donburi.Entry) (x, y, width, height float64, ok bool) {
	if !entry.HasComponent(ecs.Position) || !entry.HasComponent(ecs.Image) {
		return 0, 0, 0, 0, false
	}

	position := ecs.Position.Get(entry)
	image := *ecs.Image.Get(entry)
	if image == nil {
		return 0, 0, 0, 0, false
	}

	bounds := image.Bounds()
	return position.X, position.Y, float64(bounds.Dx()), float64(bounds.Dy()), true
}

func pointInRect(point ecs.Point, x, y, width, height float64) bool {
	return point.X >= x &&
		point.X < x+width &&
		point.Y >= y &&
		point.Y < y+height
}

func rectsOverlap(aX, aY, aWidth, aHeight, bX, bY, bWidth, bHeight float64) bool {
	return aX < bX+bWidth &&
		aX+aWidth > bX &&
		aY < bY+bHeight &&
		aY+aHeight > bY
}
