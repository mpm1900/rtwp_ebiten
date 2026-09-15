package entities

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/tiles"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
)

// CreateTileMap loads the map file, registers the TileMap component on a new
// singleton entity, and pre-renders all non-empty tiles to a single image that
// the tile renderer draws each frame. Returns the pre-rendered image.
func CreateTileMap(ecs *ecs.ECS, path string, offsetX, offsetY int) *ebiten.Image {
	tileMap, err := tiles.LoadMap(path)
	if err != nil {
		panic(err)
	}

	// Pre-render the tile layer to a world-sized image so we only pay one
	// DrawImage per frame regardless of map size. Only origin cells carry a
	// tile character; multi-cell tiles are drawn once from their top-left.
	canvas := ebiten.NewImage(components.WORLD_WIDTH, components.WORLD_HEIGHT)
	ts := tiles.TileSize
	for y := range tileMap.Height {
		for x := range tileMap.Width {
			c := tileMap.Cell(x, y)
			if tiles.IsEmpty(c) {
				continue
			}
			tile, ok := tiles.Get(c)
			if !ok {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*ts), float64(y*ts))
			canvas.DrawImage(tile.Image, op)
		}
	}

	entity := ecs.World.Create(components.TileMap)
	entry := ecs.World.Entry(entity)
	components.SetTileMap(entry, components.TileMapData{
		Map:     tileMap,
		OffsetX: offsetX,
		OffsetY: offsetY,
	})

	return canvas
}
