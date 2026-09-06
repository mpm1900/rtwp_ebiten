package components

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

const (
	SCREEN_HEIGHT = 720
	SCREEN_WIDTH  = 1280

	WORLD_BORDER = 100
	WORLD_WIDTH  = 3000
	WORLD_HEIGHT = 3000

	MINIMAP_SIZE    = 150
	MINIMAP_PADDING = 12
	MINIMAP_BORDER  = 2
)

type PlayerData struct {
	SelectedAction *Action
}

func WorldRect() (minX, minY, maxX, maxY float64) {
	return WORLD_BORDER, WORLD_BORDER, WORLD_BORDER + WORLD_WIDTH, WORLD_BORDER + WORLD_HEIGHT
}

func IsInWorld(pos math.Vec2) bool {
	minX, minY, maxX, maxY := WorldRect()
	return pos.X >= minX && pos.X <= maxX && pos.Y >= minY && pos.Y <= maxY
}

func MinimapRect() image.Rectangle {
	windowW, windowH := ebiten.WindowSize()
	mapX := windowW - MINIMAP_SIZE - MINIMAP_PADDING
	mapY := windowH - MINIMAP_SIZE - MINIMAP_PADDING
	return image.Rect(mapX, mapY, mapX+MINIMAP_SIZE, mapY+MINIMAP_SIZE)
}

func IsOverMinimap(point math.Vec2, mapRect image.Rectangle) bool {
	if point.X < float64(mapRect.Min.X) || point.Y < float64(mapRect.Min.Y) || point.X >= float64(mapRect.Max.X) || point.Y >= float64(mapRect.Max.Y) {
		return false
	}

	return true
}

func MinimapWorldPoint(point math.Vec2) (math.Vec2, bool) {
	mapRect := MinimapRect()
	if !IsOverMinimap(point, mapRect) {
		return math.Vec2{}, false
	}

	worldMinX, worldMinY, worldMaxX, worldMaxY := WorldRect()
	worldW := worldMaxX - worldMinX
	worldH := worldMaxY - worldMinY

	relX := (point.X - float64(mapRect.Min.X)) / float64(mapRect.Dx())
	relY := (point.Y - float64(mapRect.Min.Y)) / float64(mapRect.Dy())
	return math.NewVec2(worldMinX+relX*worldW, worldMinY+relY*worldH), true
}

func ClampWorldPosition(pos math.Vec2) math.Vec2 {
	minX, minY, maxX, maxY := WorldRect()
	return math.NewVec2(
		min(maxX, max(minX, pos.X)),
		min(maxY, max(minY, pos.Y)),
	)
}

func NewPlayerData() PlayerData {
	return PlayerData{
		SelectedAction: nil,
	}
}

var Player = donburi.NewComponentType[PlayerData]()

func GetPlayerEntity(world donburi.World) donburi.Entity {
	entry := Player.MustFirst(world)
	return entry.Entity()
}
func GetPlayer(world donburi.World) *PlayerData {
	entry := Player.MustFirst(world)
	return Player.Get(entry)
}
