package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type CursorData struct {
	DragStart *math.Vec2
	DragEnd   *math.Vec2
}

func (c *CursorData) StartDrag(point math.Vec2) {
	c.DragStart = &point
	c.DragEnd = nil
}

func (c *CursorData) UpdateDrag(point math.Vec2) {
	if c.DragStart == nil {
		return
	}

	c.DragEnd = &point
}
func (c *CursorData) ClearDrag() {
	c.DragStart = nil
	c.DragEnd = nil
}
func (c *CursorData) DragDistance() float64 {
	if c.DragStart == nil || c.DragEnd == nil {
		return 0
	}

	return c.DragStart.Distance(*c.DragEnd)
}
func (c *CursorData) DragRect() image.Rectangle {
	if c.DragStart == nil || c.DragEnd == nil {
		return image.Rectangle{}
	}

	return util.ToRect(*c.DragStart, *c.DragEnd)
}

var Cursor = donburi.NewComponentType[CursorData]()

func WithCursor(entry *donburi.Entry, data CursorData) {
	entry.AddComponent(Cursor)
	Cursor.SetValue(entry, data)
}

func GetCursor(world donburi.World) *CursorData {
	entry := Cursor.MustFirst(world)
	return Cursor.Get(entry)
}
