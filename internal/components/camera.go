package components

import (
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

const (
	MIN_CAMERA_ZOOM = 0.7
	MAX_CAMERA_ZOOM = 2.0
)

type CameraData struct {
	Camera     *camera.Camera
	CameraDrag *math.Vec2
}

func newCamera() *camera.Camera {
	return camera.NewCamera(SCREEN_WIDTH, SCREEN_HEIGHT, float64(WORLD_BORDER+SCREEN_WIDTH/2), float64(WORLD_BORDER+SCREEN_HEIGHT/2), 0, 1)
}
func NewCameraData() CameraData {
	return CameraData{
		Camera:     newCamera(),
		CameraDrag: nil,
	}
}

func (c *CameraData) ClampCameraPosition() {
	if c.Camera == nil {
		return
	}

	scale := c.Camera.Scale
	if scale <= 0 {
		scale = 1.0
	}

	halfWidth := float64(c.Camera.Width) / 2.0 / scale
	halfHeight := float64(c.Camera.Height) / 2.0 / scale
	minX := halfWidth
	minY := halfHeight
	maxX := WORLD_BORDER*2 + WORLD_WIDTH - halfWidth
	maxY := WORLD_BORDER*2 + WORLD_HEIGHT - halfHeight

	if minX > maxX {
		c.Camera.X = float64(WORLD_BORDER*2+WORLD_WIDTH) / 2.0
	} else {
		c.Camera.X = min(maxX, max(minX, c.Camera.X))
	}

	if minY > maxY {
		c.Camera.Y = float64(WORLD_BORDER*2+WORLD_HEIGHT) / 2.0
	} else {
		c.Camera.Y = min(maxY, max(minY, c.Camera.Y))
	}
}

func (c *CameraData) StartCameraDrag(point math.Vec2) {
	c.CameraDrag = &point
}

func (c *CameraData) UpdateCameraDrag(point math.Vec2) (math.Vec2, bool) {
	if c.CameraDrag == nil {
		c.StartCameraDrag(point)
		return math.Vec2{}, false
	}

	delta := point.Sub(*c.CameraDrag).MulScalar(2)
	c.CameraDrag = &point
	if delta.IsZero() {
		return math.Vec2{}, false
	}

	return delta, true
}

func (c *CameraData) ClearCameraDrag() {
	c.CameraDrag = nil
}

func (c *CameraData) ScreenToWorld(point math.Vec2) math.Vec2 {
	if c.Camera == nil {
		return point
	}

	x, y := c.Camera.GetWorldCoords(point.X, point.Y)
	return math.NewVec2(x, y)
}

var Camera = donburi.NewComponentType[CameraData]()

func WithCamera(entry *donburi.Entry, data CameraData) {
	entry.AddComponent(Camera)
	Camera.SetValue(entry, data)
}

func GetCamera(world donburi.World) *CameraData {
	entry := Camera.MustFirst(world)
	return Camera.Get(entry)
}
