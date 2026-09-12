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
	tile_size := components.TILE_SIZE

	// Pre-cache all sub-image tiles from the tileset to avoid 10,000 sub-image allocation calls.
	// This dramatically reduces garbage collection overhead and optimizes startup times.
	tile_images := make([]*ebiten.Image, assets.GrassTileCount)
	for i := range assets.GrassTileCount {
		col := i % assets.GrassTilesetCols
		row := i / assets.GrassTilesetCols
		tile_rect := image.Rect(
			col*tile_size,
			row*tile_size,
			(col+1)*tile_size,
			(row+1)*tile_size,
		)
		tile_images[i] = tileset.SubImage(tile_rect).(*ebiten.Image)
	}

	// Weights for the grass tiles (0 to 7). Higher weights increase probability of selection.
	// These values match the previous geometric distribution exactly:
	// Index 0 (plain grass): 50% chance. Indices 1-7 (detail tiles): geometric half chances.
	tile_weights := []int{128, 64, 32, 16, 8, 4, 2, 2}
	total_weight := 0
	for _, w := range tile_weights {
		total_weight += w
	}

	for y := range components.WORLD_TILE_COUNT {
		for x := range components.WORLD_TILE_COUNT {
			rnd_weight := rand.IntN(total_weight)
			tile_index := 0
			sum_weight := 0
			for idx, w := range tile_weights {
				sum_weight += w
				if rnd_weight < sum_weight {
					tile_index = idx
					break
				}
			}

			tile := tile_images[tile_index]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tile_size), float64(y*tile_size))
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
