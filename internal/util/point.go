package util

import (
	"image"
	"math"

	dmath "github.com/yohamta/donburi/features/math"
)

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
