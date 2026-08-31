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

func renderOutlinedSprite(screen *ebiten.Image, base *ebiten.Image, transformOptions ebiten.DrawImageOptions, outline color.Color) {
	screen.DrawImage(base, &transformOptions)

	bounds := base.Bounds()
	outlineImage := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			baseColor := color.NRGBAModel.Convert(base.At(x, y)).(color.NRGBA)
			if baseColor.A == 0 {
				continue
			}

			isEdge := false
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}

					nx := x + dx
					ny := y + dy
					if nx < bounds.Min.X || nx >= bounds.Max.X || ny < bounds.Min.Y || ny >= bounds.Max.Y {
						isEdge = true
						break
					}

					neighbor := color.NRGBAModel.Convert(base.At(nx, ny)).(color.NRGBA)
					if neighbor.A == 0 {
						isEdge = true
						break
					}
				}
				if isEdge {
					break
				}
			}

			if isEdge {
				outlineImage.Set(x-bounds.Min.X, y-bounds.Min.Y, outline)
			}
		}
	}

	outlineOptions := transformOptions
	screen.DrawImage(outlineImage, &outlineOptions)
}

func RenderActors(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)
	player := components.GetPlayerEntity(ecs.World)

	for entry := range renderActorsQuery.Iter(ecs.World) {
		trans := transform.Transform.Get(entry)
		image := *components.Image.Get(entry)
		options := ebiten.DrawImageOptions{}

		centerScale := components.CenterScale(*trans)
		options.GeoM.Translate(centerScale.X, centerScale.Y)
		options.GeoM.Rotate(dmath.ToRadians(trans.LocalRotation))

		center := components.CenterTrans(*trans)
		centerPoint := view.Point(center)
		options.GeoM.Translate(centerPoint.X, centerPoint.Y)

		outlineColor := assets.ColorEnemy
		actor := components.Actor.Get(entry)

		if actor.Player == player {
			outlineColor = assets.ColorPlayer
		}

		if entry.HasComponent(components.Selected) {
			outlineColor = assets.ColorSelected
		}

		renderOutlinedSprite(screen, image, options, outlineColor)
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
			fill.Fill(assets.ColorHpFull)
			fillOptions := &ebiten.DrawImageOptions{}
			fillOptions.GeoM.Translate(barX+1, barY+1)
			screen.DrawImage(fill, fillOptions)
		}
	}
}
