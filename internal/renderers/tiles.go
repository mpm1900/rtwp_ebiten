package renderers

import (
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
)

// tileLayer is the pre-rendered tile map image set at game startup.
var tileLayer *ebiten.Image

// SetTileLayer stores the pre-rendered tile map image for the renderer to draw.
func SetTileLayer(img *ebiten.Image) {
	tileLayer = img
}

// RenderTiles draws the pre-rendered tile map over the grass background. Tiles
// render above grass but below effects/actors so actors walk on top of them.
func RenderTiles(ecs *ecs.ECS, screen *ebiten.Image) {
	if tileLayer == nil {
		return
	}
	view := newCameraView(ecs)
	worldRect := components.WorldRect()

	op := &ebiten.DrawImageOptions{}
	origin := view.Point(dmath.NewVec2(float64(worldRect.Min.X), float64(worldRect.Min.Y)))
	op.GeoM.Translate(origin.X, origin.Y)
	screen.DrawImage(tileLayer, op)
}
