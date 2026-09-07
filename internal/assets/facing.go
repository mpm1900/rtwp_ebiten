package assets

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ActorFacingKey struct {
	Flip  bool
	Angle int
}

var actorFacingSprites map[ActorFacingKey]*ebiten.Image

func initActorFacingSprites(source *ebiten.Image) {
	actorFacingSprites = map[ActorFacingKey]*ebiten.Image{}
	angles := []int{-90, -45, 0, 45, 90}

	for _, flip := range []bool{false, true} {
		for _, angle := range angles {
			key := ActorFacingKey{Flip: flip, Angle: angle}
			actorFacingSprites[key] = bakeFacingSprite(source, flip, float64(angle))
		}
	}
}

func ActorFacingSprite(rotation float64) *ebiten.Image {
	return actorFacingSprites[facingKey(rotation)]
}

func facingKey(rotation float64) ActorFacingKey {
	flip := false
	if math.Abs(rotation) > 90 {
		flip = true
		if rotation > 0 {
			rotation -= 180
		} else {
			rotation += 180
		}
	}

	angle := int(math.Round(rotation/45)) * 45
	angle = max(-90, min(90, angle))

	return ActorFacingKey{Flip: flip, Angle: angle}
}

func bakeFacingSprite(source *ebiten.Image, flip bool, angle_deg float64) *ebiten.Image {
	bounds := source.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	canvas := int(math.Ceil(float64(width+height) / math.Sqrt2))

	dst := ebiten.NewImage(canvas, canvas)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(width)/2, -float64(height)/2)
	if flip {
		op.GeoM.Scale(-1, 1)
	}
	op.GeoM.Rotate(angle_deg * math.Pi / 180)
	op.GeoM.Translate(float64(canvas)/2, float64(canvas)/2)
	dst.DrawImage(source, op)

	return dst
}
