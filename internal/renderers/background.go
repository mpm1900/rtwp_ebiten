package renderers

import (
	"image/color"
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
	worldBackground.Fill(color.RGBA{R: 0x22, G: 0x1a, B: 0x2d, A: 0xff})
	return worldBackground
}

func RenderBackground(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)
	worldMinX, worldMinY, _, _ := components.WorldRect()

	screen.Fill(color.Black)

	playableOp := &ebiten.DrawImageOptions{}
	playableOrigin := view.Point(dmath.NewVec2(worldMinX, worldMinY))
	playableOp.GeoM.Translate(playableOrigin.X, playableOrigin.Y)
	screen.DrawImage(getWorldBackground(), playableOp)
}
