package components

import (
	"image"
	"rtwp_ebitengine/internal/tiles"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
)

// TileMapData holds the loaded tile map and its offset within world space.
type TileMapData struct {
	Map     *tiles.Map
	OffsetX int // world pixel offset of grid cell (0,0)'s top-left
	OffsetY int
}

var TileMap = donburi.NewComponentType[TileMapData]()
var TileMapQuery = donburi.NewQuery(filter.Contains(TileMap))

// SetTileMap stores the tile map on an entry.
func SetTileMap(entry *donburi.Entry, data TileMapData) {
	entry.AddComponent(TileMap)
	TileMap.SetValue(entry, data)
}

// GetTileMap returns the singleton tile map, or nil if none exists.
func GetTileMap(world donburi.World) *TileMapData {
	entry, ok := TileMapQuery.First(world)
	if !ok {
		return nil
	}
	return TileMap.Get(entry)
}

// CollidesWithTile reports whether the rect of `entry` placed at `position`
// overlaps any solid edge of any tile in the world tile map.
//
// Collision model: every solid side of a tile is a thin wall segment on that
// edge (tiles.WallThickness px). For multi-cell tiles the solid sides apply to
// the outer edges of the whole tile. Walls are never a full tile thick — a tile
// with all 4 sides solid is a sealed perimeter of thin walls; its interior is
// unreachable because every edge blocks. A partial tile (e.g. stairs with
// left/right solid) is walkable through its open sides.
func CollidesWithTile(world donburi.World, entry *donburi.Entry, position dmath.Vec2) bool {
	tm := GetTileMap(world)
	if tm == nil || tm.Map == nil {
		return false
	}
	bounds, ok := RectAt(entry, position)
	if !ok {
		return false
	}
	return tm.RectHitsTile(bounds)
}

// isFullSolid reports whether every side of the tile is solid. Used by the
// pathfinder to treat fully-sealed tiles as impassable cells.
func isFullSolid(solid [4]bool) bool {
	return solid[tiles.SideTop] && solid[tiles.SideRight] &&
		solid[tiles.SideBottom] && solid[tiles.SideLeft]
}

// tileWorldRect returns the world-space rect of a tile whose origin is at grid
// cell (ox, oy) and spans w x h cells.
func (tm *TileMapData) tileWorldRect(ox, oy, w, h int) image.Rectangle {
	ts := tiles.TileSize
	return image.Rect(
		tm.OffsetX+ox*ts, tm.OffsetY+oy*ts,
		tm.OffsetX+(ox+w)*ts, tm.OffsetY+(oy+h)*ts,
	)
}

// RectHitsTile reports whether a world-space rect overlaps any solid tile edge.
func (tm *TileMapData) RectHitsTile(bounds image.Rectangle) bool {
	ts := tiles.TileSize
	minTX := (bounds.Min.X - tm.OffsetX) / ts
	minTY := (bounds.Min.Y - tm.OffsetY) / ts
	maxTX := (bounds.Max.X - 1 - tm.OffsetX) / ts
	maxTY := (bounds.Max.Y - 1 - tm.OffsetY) / ts

	visited := map[[2]int]bool{}
	for ty := minTY; ty <= maxTY; ty++ {
		for tx := minTX; tx <= maxTX; tx++ {
			tile, ox, oy, ok := tm.tileAtOrigin(tx, ty)
			if !ok {
				continue
			}
			key := [2]int{ox, oy}
			if visited[key] {
				continue
			}
			visited[key] = true
			rect := tm.tileWorldRect(ox, oy, tile.W, tile.H)
			if rectHitsTileEdges(bounds, rect, tile.Solid) {
				return true
			}
		}
	}
	return false
}

// ForEachSolidArea yields every solid area a pathfinder should treat as an
// obstacle: fully-sealed tiles contribute their whole rect, while partial
// tiles (e.g. stairs) contribute only their thin solid edge segments. This lets
// the pathfinder route through the open middle of a wide stair while avoiding
// its railings.
func (tm *TileMapData) ForEachSolidArea(yield func(image.Rectangle)) {
	if tm == nil || tm.Map == nil {
		return
	}
	for ty := range tm.Map.Height {
		for tx := range tm.Map.Width {
			c := tm.Map.Cells[ty*tm.Map.Width+tx]
			if tiles.IsEmpty(c) {
				continue
			}
			tile, ok := tiles.Get(c)
			if !ok {
				continue
			}
			rect := tm.tileWorldRect(tx, ty, tile.W, tile.H)
			if isFullSolid(tile.Solid) {
				yield(rect)
				continue
			}
			thick := tiles.WallThickness
			if tile.Solid[tiles.SideTop] {
				yield(image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y+thick))
			}
			if tile.Solid[tiles.SideBottom] {
				yield(image.Rect(rect.Min.X, rect.Max.Y-thick, rect.Max.X, rect.Max.Y))
			}
			if tile.Solid[tiles.SideLeft] {
				yield(image.Rect(rect.Min.X, rect.Min.Y, rect.Min.X+thick, rect.Max.Y))
			}
			if tile.Solid[tiles.SideRight] {
				yield(image.Rect(rect.Max.X-thick, rect.Min.Y, rect.Max.X, rect.Max.Y))
			}
		}
	}
}

// tileAtOrigin returns the tile occupying cell (x, y) along with its origin
// cell coordinates, resolving multi-cell tiles.
func (tm *TileMapData) tileAtOrigin(x, y int) (tiles.Tile, int, int, bool) {
	if x < 0 || y < 0 || x >= tm.Map.Width || y >= tm.Map.Height {
		return tiles.Tile{}, 0, 0, false
	}
	idx := y*tm.Map.Width + x
	ox := tm.Map.OriginX[idx]
	oy := tm.Map.OriginY[idx]
	if ox < 0 {
		return tiles.Tile{}, 0, 0, false
	}
	c := tm.Map.Cells[oy*tm.Map.Width+ox]
	if tiles.IsEmpty(c) {
		return tiles.Tile{}, 0, 0, false
	}
	tile, ok := tiles.Get(c)
	if !ok {
		return tiles.Tile{}, 0, 0, false
	}
	return tile, ox, oy, true
}

// rectHitsTileEdges checks `bounds` against each solid edge segment of a tile
// whose world rect is `rect`. Each solid side is a thin wall on that edge.
func rectHitsTileEdges(bounds, rect image.Rectangle, solid [4]bool) bool {
	thick := tiles.WallThickness
	if solid[tiles.SideTop] {
		seg := image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y+thick)
		if bounds.Overlaps(seg) {
			return true
		}
	}
	if solid[tiles.SideBottom] {
		seg := image.Rect(rect.Min.X, rect.Max.Y-thick, rect.Max.X, rect.Max.Y)
		if bounds.Overlaps(seg) {
			return true
		}
	}
	if solid[tiles.SideLeft] {
		seg := image.Rect(rect.Min.X, rect.Min.Y, rect.Min.X+thick, rect.Max.Y)
		if bounds.Overlaps(seg) {
			return true
		}
	}
	if solid[tiles.SideRight] {
		seg := image.Rect(rect.Max.X-thick, rect.Min.Y, rect.Max.X, rect.Max.Y)
		if bounds.Overlaps(seg) {
			return true
		}
	}
	return false
}
