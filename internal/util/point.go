package util

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	dmath "github.com/yohamta/donburi/features/math"
)

func NewVec2(x int, y int) dmath.Vec2 {
	return dmath.NewVec2(float64(x), float64(y))
}

func ToPoint(v dmath.Vec2) image.Point {
	return image.Pt(
		int(math.Floor(v.X)),
		int(math.Floor(v.Y)),
	)
}

func ToRect(start, end dmath.Vec2) image.Rectangle {
	min := image.Pt(
		int(math.Floor(math.Min(start.X, end.X))),
		int(math.Floor(math.Min(start.Y, end.Y))),
	)
	max := image.Pt(
		int(math.Ceil(math.Max(start.X, end.X))),
		int(math.Ceil(math.Max(start.Y, end.Y))),
	)
	return image.Rectangle{
		Min: min,
		Max: max,
	}
}

func DrawPoints(screen *ebiten.Image, start, end dmath.Vec2, strokeWidth float32, strokeColor color.Color) {
	vector.StrokeLine(
		screen,
		float32(start.X),
		float32(start.Y),
		float32(end.X),
		float32(end.Y),
		strokeWidth,
		strokeColor,
		true,
	)
}

func CursorPoint() dmath.Vec2 {
	x, y := ebiten.CursorPosition()
	return NewVec2(x, y)
}
