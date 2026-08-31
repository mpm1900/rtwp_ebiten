package renderers

import (
	"image/color"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
)

func RenderBackground(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)
	worldMinX, worldMinY, _, _ := components.WorldRect()

	screen.Fill(color.Black)

	playable := ebiten.NewImage(int(components.WorldWidth), int(components.WorldHeight))
	playable.Fill(color.RGBA{R: 0x22, G: 0x1a, B: 0x2d, A: 0xff})

	playableOp := &ebiten.DrawImageOptions{}
	playableOrigin := dmath.NewVec2(worldMinX, worldMinY)
	playableOrigin = view.Point(playableOrigin)
	playableOp.GeoM.Translate(playableOrigin.X, playableOrigin.Y)
	screen.DrawImage(playable, playableOp)
}
