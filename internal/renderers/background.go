package renderers

import (
	"image"
	"image/color"
	"math/rand/v2"
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
)

var worldBackground *ebiten.Image

func getWorldBackground() *ebiten.Image {
	if worldBackground != nil {
		return worldBackground
	}

	worldBackground = ebiten.NewImage(int(components.WORLD_WIDTH), int(components.WORLD_HEIGHT))
	tileset := assets.GrassTilesetImage
	tileSize := components.TILE_SIZE

	for y := range components.WORLD_TILE_COUNT {
		for x := range components.WORLD_TILE_COUNT {
			tileIndex := 0
			if rand.IntN(2) > 0 {
				tileIndex = 1
				if rand.IntN(2) > 0 {
					tileIndex = 2
					if rand.IntN(2) > 0 {
						tileIndex = 3
						if rand.IntN(2) > 0 {
							tileIndex = 4
							if rand.IntN(2) > 0 {
								tileIndex = 5
								if rand.IntN(2) > 0 {
									tileIndex = 6
									if rand.IntN(2) > 0 {
										tileIndex = 7
									}
								}
							}
						}
					}
				}
			}

			col := tileIndex % assets.GrassTilesetCols
			row := tileIndex / assets.GrassTilesetCols
			tileRect := image.Rect(
				col*tileSize,
				row*tileSize,
				(col+1)*tileSize,
				(row+1)*tileSize,
			)
			tile := tileset.SubImage(tileRect).(*ebiten.Image)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			worldBackground.DrawImage(tile, op)
		}
	}

	return worldBackground
}

func RenderBackground(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)
	worldRect := components.WorldRect()

	screen.Fill(color.Black)

	playableOp := &ebiten.DrawImageOptions{}
	playableOrigin := view.Point(dmath.NewVec2(float64(worldRect.Min.X), float64(worldRect.Min.Y)))
	playableOp.GeoM.Translate(playableOrigin.X, playableOrigin.Y)
	screen.DrawImage(getWorldBackground(), playableOp)
}
