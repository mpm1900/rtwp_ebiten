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
var CursorMoveImage *ebiten.Image

var ActorImage *ebiten.Image
var BlueSquareImage *ebiten.Image
var GreenSquareImage *ebiten.Image
var YellowSquareImage *ebiten.Image

var YolkFontSource *text.GoTextFaceSource

func MustLoadAssets() {
	var err error
	CursorPointerImage, _, _ = ebitenutil.NewImageFromFile("assets/images/cursor-pointer.png")
	CursorMoveImage, _, err = ebitenutil.NewImageFromFile("assets/images/cursor-move.png")
	ActorImage = ebiten.NewImage(24, 24)
	ActorImage.Fill(ColorActor)
	BlueSquareImage = ebiten.NewImage(24, 24)
	BlueSquareImage.Fill(color.RGBA{0, 0, 0xff, 0xff})
	GreenSquareImage = ebiten.NewImage(24, 24)
	GreenSquareImage.Fill(color.RGBA{0, 0xff, 0, 0xff})
	YellowSquareImage = ebiten.NewImage(12, 12)
	YellowSquareImage.Fill(color.RGBA{0xff, 0xff, 0, 0xff})

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Yolk6TTF))
	if err != nil {
		log.Fatal(err)
	}
	YolkFontSource = s
}
