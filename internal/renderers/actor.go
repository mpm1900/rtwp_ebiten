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
	HEALTH_BAR_WIDTH  float32 = 28
)

var renderActorsQuery = donburi.NewQuery(filter.Contains(components.Actor, components.Image))
var outlineCache = map[*ebiten.Image]map[uint32]*ebiten.Image{}

func outlineColorKey(outline color.Color, thickness int) uint32 {
	rgba := color.RGBAModel.Convert(outline).(color.RGBA)
	return uint32(rgba.R)<<24 | uint32(rgba.G)<<16 | uint32(rgba.B)<<8 | uint32(rgba.A)<<0 | uint32(thickness)
}

func buildOutlineImage(base *ebiten.Image, outline color.Color, thickness int) *ebiten.Image {
	if thickness <= 0 {
		return ebiten.NewImage(base.Bounds().Dx(), base.Bounds().Dy())
	}

	bounds := base.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	outlineImage := ebiten.NewImage(width+thickness*2, height+thickness*2)
	rgba := color.RGBAModel.Convert(outline).(color.RGBA)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			baseColor := color.NRGBAModel.Convert(base.At(x, y)).(color.NRGBA)
			if baseColor.A == 0 {
				continue
			}

			localX := x - bounds.Min.X
			localY := y - bounds.Min.Y
			for dx := -thickness; dx <= thickness; dx++ {
				for dy := -thickness; dy <= thickness; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}

					nx := x + dx
					ny := y + dy
					if nx < bounds.Min.X || nx >= bounds.Max.X || ny < bounds.Min.Y || ny >= bounds.Max.Y {
						outlineImage.Set(localX+thickness+dx, localY+thickness+dy, rgba)
						continue
					}

					neighbor := color.NRGBAModel.Convert(base.At(nx, ny)).(color.NRGBA)
					if neighbor.A == 0 {
						outlineImage.Set(localX+thickness+dx, localY+thickness+dy, rgba)
					}
				}
			}
		}
	}

	return outlineImage
}

func getOutlineImage(base *ebiten.Image, outline color.Color, thickness int) *ebiten.Image {
	key := outlineColorKey(outline, thickness)
	if _, ok := outlineCache[base]; !ok {
		outlineCache[base] = map[uint32]*ebiten.Image{}
	}
	if cached, ok := outlineCache[base][key]; ok {
		return cached
	}

	cached := buildOutlineImage(base, outline, thickness)
	outlineCache[base][key] = cached
	return cached
}

func renderOutlinedSprite(screen *ebiten.Image, base *ebiten.Image, transformOptions ebiten.DrawImageOptions, outline color.Color, thickness int) {
	screen.DrawImage(base, &transformOptions)
	outlineImage := getOutlineImage(base, outline, thickness)
	outlineOptions := ebiten.DrawImageOptions{}
	outlineOptions.GeoM.Translate(float64(-thickness), float64(-thickness))
	outlineOptions.GeoM.Concat(transformOptions.GeoM)
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
		thickness := 1

		if actor.Player == player {
			outlineColor = assets.ColorPlayer
		}

		if entry.HasComponent(components.Selected) {
			outlineColor = assets.ColorSelected
			thickness = 2
		}

		renderOutlinedSprite(screen, image, options, outlineColor, thickness)
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
		barX := actorCenter.X - 14
		barY := actorCenter.Y - 18

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
