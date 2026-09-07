package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/features/math"
)

const (
	SCREEN_HEIGHT = 720
	SCREEN_WIDTH  = 1280

	TILE_SIZE        = 32
	WORLD_TILE_COUNT = 100

	WORLD_BORDER = 100
	WORLD_WIDTH  = WORLD_TILE_COUNT * TILE_SIZE
	WORLD_HEIGHT = WORLD_TILE_COUNT * TILE_SIZE

	MINIMAP_SIZE    = 150
	MINIMAP_PADDING = 12
	MINIMAP_BORDER  = 2
)

var worldRect = image.Rect(
	WORLD_BORDER, WORLD_BORDER,
	WORLD_BORDER+WORLD_WIDTH, WORLD_BORDER+WORLD_HEIGHT,
)

func WorldRect() image.Rectangle {
	return worldRect
}

func IsInWorld(pos math.Vec2) bool {
	p := util.ToPoint(pos)
	return p.X >= worldRect.Min.X && p.X < worldRect.Max.X &&
		p.Y >= worldRect.Min.Y && p.Y < worldRect.Max.Y
}

func MinimapRect() image.Rectangle {
	windowW, windowH := ebiten.WindowSize()
	mapX := windowW - MINIMAP_SIZE - MINIMAP_PADDING
	mapY := windowH - MINIMAP_SIZE - MINIMAP_PADDING
	return image.Rect(mapX, mapY, mapX+MINIMAP_SIZE, mapY+MINIMAP_SIZE)
}

func MinimapWorldPoint(point math.Vec2) (math.Vec2, bool) {
	mapRect := MinimapRect()
	p := util.ToPoint(point)
	if p.X < mapRect.Min.X || p.Y < mapRect.Min.Y || p.X >= mapRect.Max.X || p.Y >= mapRect.Max.Y {
		return math.Vec2{}, false
	}

	relX := (point.X - float64(mapRect.Min.X)) / float64(mapRect.Dx())
	relY := (point.Y - float64(mapRect.Min.Y)) / float64(mapRect.Dy())
	return math.NewVec2(
		float64(worldRect.Min.X)+relX*float64(worldRect.Dx()),
		float64(worldRect.Min.Y)+relY*float64(worldRect.Dy()),
	), true
}

func ClampWorldPosition(pos math.Vec2) math.Vec2 {
	return math.NewVec2(
		min(float64(worldRect.Max.X), max(float64(worldRect.Min.X), pos.X)),
		min(float64(worldRect.Max.Y), max(float64(worldRect.Min.Y), pos.Y)),
	)
}
