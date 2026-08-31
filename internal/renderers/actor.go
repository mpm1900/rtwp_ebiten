package renderers

import (
	"image/color"
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

const (
	HEALTH_BAR_HEIGHT float32 = 6
	HEALTH_BAR_WIDTH  float32 = 24
)

var renderActorsQuery = donburi.NewQuery(filter.Contains(components.Actor, components.Image))
var outlineCache = map[*ebiten.Image]map[uint32]*ebiten.Image{}

func outlineColorKey(outline color.Color) uint32 {
	rgba := color.RGBAModel.Convert(outline).(color.RGBA)
	return uint32(rgba.R)<<24 | uint32(rgba.G)<<16 | uint32(rgba.B)<<8 | uint32(rgba.A)
}

func buildOutlineImage(base *ebiten.Image, outline color.Color) *ebiten.Image {
	bounds := base.Bounds()
	outlineImage := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	rgba := color.RGBAModel.Convert(outline).(color.RGBA)

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
				outlineImage.Set(x-bounds.Min.X, y-bounds.Min.Y, color.RGBA{R: rgba.R, G: rgba.G, B: rgba.B, A: rgba.A})
			}
		}
	}

	return outlineImage
}

func getOutlineImage(base *ebiten.Image, outline color.Color) *ebiten.Image {
	key := outlineColorKey(outline)
	if _, ok := outlineCache[base]; !ok {
		outlineCache[base] = map[uint32]*ebiten.Image{}
	}
	if cached, ok := outlineCache[base][key]; ok {
		return cached
	}

	cached := buildOutlineImage(base, outline)
	outlineCache[base][key] = cached
	return cached
}

func renderOutlinedSprite(screen *ebiten.Image, base *ebiten.Image, transformOptions ebiten.DrawImageOptions, outline color.Color) {
	screen.DrawImage(base, &transformOptions)
	outlineImage := getOutlineImage(base, outline)
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
		if health <= 0 {
			continue
		}

		percent := float32((health - damage) / health)
		if percent < 0 {
			percent = 0
		}
		if percent > 1 {
			percent = 1
		}

		actorCenter := view.Point(components.CenterTrans(*trans))
		barX := actorCenter.X - 12
		barY := actorCenter.Y - 12

		vector.FillRect(
			screen,
			float32(barX),
			float32(barY),
			HEALTH_BAR_WIDTH,
			HEALTH_BAR_HEIGHT,
			color.RGBA{0, 0, 0, 0xff},
			false,
		)
		if percent > 0 {
			vector.FillRect(
				screen,
				float32(barX+1),
				float32(barY+1),
				HEALTH_BAR_WIDTH*percent-2,
				HEALTH_BAR_HEIGHT-2,
				assets.ColorHpFull,
				false,
			)
		}
	}
}
