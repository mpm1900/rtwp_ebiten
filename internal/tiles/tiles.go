package tiles

import (
	"image"
	"rtwp_ebitengine/internal/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// Side indices for the Solid array on a Tile.
const (
	SideTop    = 0
	SideRight  = 1
	SideBottom = 2
	SideLeft   = 3
)

// TileSize matches components.TILE_SIZE; duplicated here to avoid an import cycle
// between tiles and components.
const TileSize = 32

// Tile is a renderable map tile with optional per-side collision. A tile may
// span multiple cells (W x H), e.g. a 2x3 staircase sprite. Solid sides are
// treated as thin wall segments on the OUTER edges of the whole tile.
type Tile struct {
	Image *ebiten.Image
	// W, H are the tile's size in cells (1x1 by default).
	W, H int
	// Solid marks which outer edges block movement: top, right, bottom, left.
	Solid [4]bool
}

// WallThickness is how thick each solid edge segment is in pixels. Kept small
// so actors can pass through the open direction of a partially-solid tile
// (e.g. stairs with left/right railings).
const WallThickness = 2

var registry = map[byte]Tile{}

// Define registers a 1x1 tile under the given single-byte id by slicing a cell
// out of the supplied tileset.
func Define(id byte, tileset *ebiten.Image, cellX, cellY int, solid [4]bool) Tile {
	return DefineMulti(id, tileset, cellX, cellY, 1, 1, solid)
}

// DefineMulti registers a tile spanning w x h cells, with its top-left cell at
// (cellX, cellY) in the tileset. The solid sides apply to the outer edges of
// the whole w x h region.
func DefineMulti(id byte, tileset *ebiten.Image, cellX, cellY, w, h int, solid [4]bool) Tile {
	rect := image.Rect(
		cellX*TileSize, cellY*TileSize,
		(cellX+w)*TileSize, (cellY+h)*TileSize,
	)
	t := Tile{
		Image: tileset.SubImage(rect).(*ebiten.Image),
		W:     w,
		H:     h,
		Solid: solid,
	}
	registry[id] = t
	return t
}

// Wall segment source rects in TX Tileset Wall.png (pixel coordinates):
//   - horizontal wall (for top/bottom edges): 7px tall, ~3 tiles wide.
//   - vertical wall with shadow (for left edges): 9px wide (7 wall + 2 shadow).
//   - vertical wall plain (for right edges): 7px wide.
var (
	hWallSrc = image.Rect(384, 32, 384+95, 32+7)
	vWallShadowSrc = image.Rect(288, 32, 288+9, 32+95)
	vWallPlainSrc = image.Rect(344, 32, 344+7, 32+95)
)

// DefineWall registers a 1x1 tile that is a thin wall on the given sides. The
// wall segment is composited onto a 32x32 cell at the chosen edges; the rest of
// the cell is transparent. Collision is enabled on exactly those sides.
func DefineWall(id byte, tileset *ebiten.Image, sides []int) Tile {
	cell := ebiten.NewImage(TileSize, TileSize)
	solid := [4]bool{}
	for _, side := range sides {
		solid[side] = true
		seg := wallSegment(tileset, side)
		op := &ebiten.DrawImageOptions{}
		switch side {
		case SideTop:
			op.GeoM.Translate(0, 0)
		case SideBottom:
			op.GeoM.Translate(0, float64(TileSize-seg.Bounds().Dy()))
		case SideLeft:
			op.GeoM.Translate(0, 0)
		case SideRight:
			op.GeoM.Translate(float64(TileSize-seg.Bounds().Dx()), 0)
		}
		cell.DrawImage(seg, op)
	}
	t := Tile{Image: cell, W: 1, H: 1, Solid: solid}
	registry[id] = t
	return t
}

// wallSegment slices a 32px-long portion of the wall segment for the given side.
func wallSegment(tileset *ebiten.Image, side int) *ebiten.Image {
	var r image.Rectangle
	switch side {
	case SideTop, SideBottom:
		// horizontal: take the first 32px of the 95px-wide segment.
		r = image.Rect(hWallSrc.Min.X, hWallSrc.Min.Y, hWallSrc.Min.X+TileSize, hWallSrc.Max.Y)
	case SideLeft:
		// vertical with shadow: 9px wide, first 32px tall.
		r = image.Rect(vWallShadowSrc.Min.X, vWallShadowSrc.Min.Y, vWallShadowSrc.Max.X, vWallShadowSrc.Min.Y+TileSize)
	case SideRight:
		// vertical plain: 7px wide, first 32px tall.
		r = image.Rect(vWallPlainSrc.Min.X, vWallPlainSrc.Min.Y, vWallPlainSrc.Max.X, vWallPlainSrc.Min.Y+TileSize)
	}
	return tileset.SubImage(r).(*ebiten.Image)
}

// Get returns the tile registered under id. ok is false if undefined.
func Get(id byte) (Tile, bool) {
	t, ok := registry[id]
	return t, ok
}

// LoadDefinitions registers all game tiles against the loaded tilesets.
// Cell coordinates point at sprites that actually contain pixels (the top-left
// 0,0 cell of these tilesets is empty).
func LoadDefinitions() {
	wall := assets.WallTilesetImage
	structs := assets.StructTilesetImage

	// Thin walls: a wall on a single edge of a 32x32 cell. The wall segment is
	// composited onto that edge; collision is enabled only on that side.
	//   n = top (north) wall    b = bottom (south) wall
	//   l = left (west) wall    r = right (east) wall
	DefineWall('n', wall, []int{SideTop})
	DefineWall('b', wall, []int{SideBottom})
	DefineWall('l', wall, []int{SideLeft})
	DefineWall('r', wall, []int{SideRight})

	// Corner walls (two edges solid) for sealing enclosures.
	DefineWall('1', wall, []int{SideTop, SideLeft})
	DefineWall('2', wall, []int{SideTop, SideRight})
	DefineWall('3', wall, []int{SideBottom, SideLeft})
	DefineWall('4', wall, []int{SideBottom, SideRight})

	// Pillar: all four sides solid (sealed thin-wall perimeter).
	DefineWall('o', structs, []int{SideTop, SideRight, SideBottom, SideLeft})

	// Stairs: three 2x3 stair sprites in TX Struct at cell columns 1, 4, 7
	// (rows 1-3). They run in the up/down direction, so the left and right
	// sides are solid (railings) and the top/bottom are open for traversal.
	DefineMulti('s', structs, 1, 1, 2, 3, [4]bool{false, true, false, true})
	DefineMulti('S', structs, 4, 1, 2, 3, [4]bool{false, true, false, true})
	DefineMulti('t', structs, 7, 1, 2, 3, [4]bool{false, true, false, true})
}
