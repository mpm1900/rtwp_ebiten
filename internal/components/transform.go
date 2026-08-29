package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func WithTransform(entry *donburi.Entry, trans transform.TransformData) {
	entry.AddComponent(transform.Transform)
	transform.Transform.SetValue(entry, trans)
}

func Rect(entry *donburi.Entry) (image.Rectangle, bool) {
	trans := transform.Transform.Get(entry)
	return RectAt(entry, trans.LocalPosition)
}

func RectAt(entry *donburi.Entry, position dmath.Vec2) (image.Rectangle, bool) {
	trans := transform.Transform.Get(entry)
	if trans.LocalScale.X <= 0 || trans.LocalScale.Y <= 0 {
		return image.Rectangle{}, false
	}

	return util.ToRect(position, position.Add(trans.LocalScale)), true
}

func Center(entry *donburi.Entry) dmath.Vec2 {
	trans := transform.Transform.Get(entry)
	return trans.LocalPosition.Add(trans.LocalScale.DivScalar(2))
}
