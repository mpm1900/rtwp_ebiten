package assets

import (
	"bytes"
	"image/color"
	"log"
	"rtwp_ebitengine/assets/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var CursorPointerImage *ebiten.Image
var CursorInvalidImage *ebiten.Image
var CursorMoveImage *ebiten.Image
var CursorAttackImage *ebiten.Image
var CursorInteractImage *ebiten.Image

var ActorImage *ebiten.Image
var BlueSquareImage *ebiten.Image
var FireImage *ebiten.Image
var GreenSquareImage *ebiten.Image
var YellowSquareImage *ebiten.Image

var YolkFontSource *text.GoTextFaceSource

var GrassTilesetImage *ebiten.Image

const (
	GrassTilesetCols = 8
	GrassTileCount   = GrassTilesetCols * GrassTilesetCols
)

func Load() {
	var err error
	CursorPointerImage, _, _ = ebitenutil.NewImageFromFile("assets/images/cursor-pointer.png")
	CursorInvalidImage, _, _ = ebitenutil.NewImageFromFile("assets/images/cursor-invalid.png")
	CursorMoveImage, _, _ = ebitenutil.NewImageFromFile("assets/images/cursor-move.png")
	CursorAttackImage, _, _ = ebitenutil.NewImageFromFile("assets/images/cursor-attack.png")
	CursorInteractImage, _, _ = ebitenutil.NewImageFromFile("assets/images/cursor-interact.png")
	cultist, _, err := ebitenutil.NewImageFromFile("assets/images/crowned-skull.png")
	if err != nil {
		log.Fatal(err)
	}
	fire, _, err := ebitenutil.NewImageFromFile("assets/images/Fire7.png")
	if err != nil {
		log.Fatal(err)
	}
	ActorImage = resizeImage(cultist, 24, 24)
	initActorFacingSprites(ActorImage)
	FireImage = fire
	BlueSquareImage = ebiten.NewImage(24, 24)
	BlueSquareImage.Fill(color.RGBA{0, 0, 0xff, 0xff})
	GreenSquareImage = ebiten.NewImage(24, 24)
	GreenSquareImage.Fill(color.RGBA{0, 0xff, 0, 0xff})
	YellowSquareImage = ebiten.NewImage(48, 48)
	YellowSquareImage.Fill(color.RGBA{0xff, 0xff, 0, 0xff})

	GrassTilesetImage, _, err = ebitenutil.NewImageFromFile("assets/images/TX Tileset Grass.png")
	if err != nil {
		log.Fatal(err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Yolk6TTF))
	if err != nil {
		log.Fatal(err)
	}
	YolkFontSource = s
}

func resizeImage(source *ebiten.Image, width, height int) *ebiten.Image {
	dst := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	bounds := source.Bounds()
	op.GeoM.Scale(
		float64(width)/float64(bounds.Dx()),
		float64(height)/float64(bounds.Dy()),
	)
	dst.DrawImage(source, op)
	return dst
}
