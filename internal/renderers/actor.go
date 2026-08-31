package renderers

import (
	"image/color"
	"math"
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var renderActorsQuery = donburi.NewQuery(filter.Contains(components.Actor, components.Image))

func RenderActors(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)

	for entry := range renderActorsQuery.Iter(ecs.World) {
		trans := transform.Transform.Get(entry)
		image := *components.Image.Get(entry)

		if entry.HasComponent(components.Selected) {
			image = assets.GreenSquareImage
		}

		options := ebiten.DrawImageOptions{}

		centerScale := components.CenterScale(*trans)
		options.GeoM.Translate(centerScale.X, centerScale.Y)

		options.GeoM.Rotate(dmath.ToRadians(trans.LocalRotation))

		center := components.CenterTrans(*trans)
		centerPoint := view.Point(center)
		options.GeoM.Translate(centerPoint.X, centerPoint.Y)

		screen.DrawImage(image, &options)
	}
}

func RenderHealthbars(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)

	for entry := range components.StatsQuery.Iter(ecs.World) {
		trans := transform.Transform.Get(entry)
		if trans == nil {
			continue
		}

		health, damage := components.GetHealth(entry)
		percent := (health - damage) / health
		if percent < 0 {
			percent = 0
		}
		if percent > 1 {
			percent = 1
		}

		actorCenter := view.Point(components.CenterTrans(*trans))
		barWidth := 24.0
		barHeight := 6.0
		barX := actorCenter.X - barWidth/2
		barY := actorCenter.Y + trans.LocalScale.DivScalar(2.0).Y

		bg := ebiten.NewImage(int(barWidth), int(barHeight))
		bg.Fill(color.RGBA{0x20, 0x20, 0x20, 0xff})
		bgOptions := &ebiten.DrawImageOptions{}
		bgOptions.GeoM.Translate(barX, barY)
		screen.DrawImage(bg, bgOptions)

		fillWidth := int(math.Max(0, barWidth*percent))
		if fillWidth > 0 {
			fill := ebiten.NewImage(fillWidth-2, int(barHeight)-2)
			fill.Fill(color.RGBA{0x4c, 0xd3, 0x66, 0xff})
			fillOptions := &ebiten.DrawImageOptions{}
			fillOptions.GeoM.Translate(barX+1, barY+1)
			screen.DrawImage(fill, fillOptions)
		}
	}
}
