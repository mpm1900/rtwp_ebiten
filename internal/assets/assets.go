package assets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var RedSquareImage *ebiten.Image
var BlueSquareImage *ebiten.Image
var GreenSquareImage *ebiten.Image
var YellowSquareImage *ebiten.Image

func MustLoadAssets() {
	RedSquareImage = ebiten.NewImage(24, 24)
	RedSquareImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
	BlueSquareImage = ebiten.NewImage(24, 24)
	BlueSquareImage.Fill(color.RGBA{0, 0, 0xff, 0xff})
	GreenSquareImage = ebiten.NewImage(24, 24)
	GreenSquareImage.Fill(color.RGBA{0, 0xff, 0, 0xff})
	YellowSquareImage = ebiten.NewImage(12, 12)
	YellowSquareImage.Fill(color.RGBA{0xff, 0xff, 0, 0xff})
}
