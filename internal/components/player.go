package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"

	camera "github.com/melonfunction/ebiten-camera"
)

const (
	SCREEN_HEIGHT = 720
	SCREEN_WIDTH  = 1280

	WORLD_BORDER = 100
	WORLD_WIDTH  = 3000
	WORLD_HEIGHT = 3000

	MINIMAP_SIZE    = 150
	MINIMAP_PADDING = 12
	MINIMAP_BORDER  = 2

	MIN_CAMERA_ZOOM = 0.7
	MAX_CAMERA_ZOOM = 2.0
)

type ActionEvent struct {
	Action *Action
	Keys   []ebiten.Key
	Loop   bool
	Point  math.Vec2
	Source donburi.Entity
}

type ActionStatus int

const (
	ActionRunning ActionStatus = iota
	ActionComplete
	ActionCanceled
)

type Action struct {
	ActionConfig
	Behavior ActionBehavior
}

type ActionConfig struct {
	Key           ebiten.Key
	Name          string
	Delay         int
	Cooldown      int
	Cursor        *ebiten.Image
	CursorInvalid *ebiten.Image
	CursorOffset  math.Vec2
}

type ActionBehavior interface {
	// Publish turns input into queued ActionEvents for selected actors.
	Publish(*Action, donburi.World, ActionEvent)
	// Start applies one-time setup for a queued ActionEvent.
	Start(*Action, donburi.World, ActionEvent)
	// Update advances an active ActionEvent and reports its lifecycle state.
	Update(*Action, donburi.World, ActionEvent) ActionStatus
	// Cancel cleans up when the action is interrupted or popped.
	Cancel(*Action, donburi.World, ActionEvent)
	// Valid checks whether input at point can use this action.
	Valid(*Action, donburi.World, math.Vec2) bool
}

func (a *Action) Publish(world donburi.World, event ActionEvent) {
	if a == nil || a.Behavior == nil {
		return
	}

	a.Behavior.Publish(a, world, event)
}

func (a *Action) Start(world donburi.World, event ActionEvent) {
	if a == nil || a.Behavior == nil {
		return
	}

	a.Behavior.Start(a, world, event)
}

func (a *Action) Update(world donburi.World, event ActionEvent) ActionStatus {
	if a == nil || a.Behavior == nil {
		return ActionComplete
	}

	return a.Behavior.Update(a, world, event)
}

func (a *Action) Cancel(world donburi.World, event ActionEvent) {
	if a == nil || a.Behavior == nil {
		return
	}

	a.Behavior.Cancel(a, world, event)
}

func (a *Action) Valid(world donburi.World, point math.Vec2) bool {
	if a == nil || a.Behavior == nil {
		return false
	}

	return a.Behavior.Valid(a, world, point)
}

type PlayerData struct {
	SelectedAction *Action
	Camera         *camera.Camera
	CameraDrag     *math.Vec2
	DragStart      *math.Vec2
	DragEnd        *math.Vec2
}

func WorldRect() (minX, minY, maxX, maxY float64) {
	return WORLD_BORDER, WORLD_BORDER, WORLD_BORDER + WORLD_WIDTH, WORLD_BORDER + WORLD_HEIGHT
}

func IsInWorld(pos math.Vec2) bool {
	minX, minY, maxX, maxY := WorldRect()
	return pos.X >= minX && pos.X <= maxX && pos.Y >= minY && pos.Y <= maxY
}

func MinimapRect() image.Rectangle {
	windowW, windowH := ebiten.WindowSize()
	mapX := windowW - MINIMAP_SIZE - MINIMAP_PADDING
	mapY := windowH - MINIMAP_SIZE - MINIMAP_PADDING
	return image.Rect(mapX, mapY, mapX+MINIMAP_SIZE, mapY+MINIMAP_SIZE)
}

func IsOverMinimap(point math.Vec2, mapRect image.Rectangle) bool {
	if point.X < float64(mapRect.Min.X) || point.Y < float64(mapRect.Min.Y) || point.X >= float64(mapRect.Max.X) || point.Y >= float64(mapRect.Max.Y) {
		return false
	}

	return true
}

func MinimapWorldPoint(point math.Vec2) (math.Vec2, bool) {
	mapRect := MinimapRect()
	if !IsOverMinimap(point, mapRect) {
		return math.Vec2{}, false
	}

	worldMinX, worldMinY, worldMaxX, worldMaxY := WorldRect()
	worldW := worldMaxX - worldMinX
	worldH := worldMaxY - worldMinY

	relX := (point.X - float64(mapRect.Min.X)) / float64(mapRect.Dx())
	relY := (point.Y - float64(mapRect.Min.Y)) / float64(mapRect.Dy())
	return math.NewVec2(worldMinX+relX*worldW, worldMinY+relY*worldH), true
}

func ClampWorldPosition(pos math.Vec2) math.Vec2 {
	minX, minY, maxX, maxY := WorldRect()
	return math.NewVec2(
		min(maxX, max(minX, pos.X)),
		min(maxY, max(minY, pos.Y)),
	)
}

func (p *PlayerData) ClampCameraPosition() {
	if p.Camera == nil {
		return
	}

	scale := p.Camera.Scale
	if scale <= 0 {
		scale = 1.0
	}

	halfWidth := float64(p.Camera.Width) / 2.0 / scale
	halfHeight := float64(p.Camera.Height) / 2.0 / scale
	minX := halfWidth
	minY := halfHeight
	maxX := WORLD_BORDER*2 + WORLD_WIDTH - halfWidth
	maxY := WORLD_BORDER*2 + WORLD_HEIGHT - halfHeight

	if minX > maxX {
		p.Camera.X = float64(WORLD_BORDER*2+WORLD_WIDTH) / 2.0
	} else {
		p.Camera.X = min(maxX, max(minX, p.Camera.X))
	}

	if minY > maxY {
		p.Camera.Y = float64(WORLD_BORDER*2+WORLD_HEIGHT) / 2.0
	} else {
		p.Camera.Y = min(maxY, max(minY, p.Camera.Y))
	}
}

func NewPlayerCamera() *camera.Camera {
	return camera.NewCamera(SCREEN_WIDTH, SCREEN_HEIGHT, float64(WORLD_BORDER+SCREEN_WIDTH/2), float64(WORLD_BORDER+SCREEN_HEIGHT/2), 0, 1)
}

func NewPlayerData() PlayerData {
	return PlayerData{
		SelectedAction: nil,
		Camera:         NewPlayerCamera(),
		CameraDrag:     nil,
		DragStart:      nil,
		DragEnd:        nil,
	}
}

func (p *PlayerData) StartDrag(point math.Vec2) {
	p.DragStart = &point
	p.DragEnd = nil
}

func (p *PlayerData) UpdateDrag(point math.Vec2) {
	if p.DragStart == nil {
		return
	}

	p.DragEnd = &point
}
func (p *PlayerData) ClearDrag() {
	p.DragStart = nil
	p.DragEnd = nil
}

func (p *PlayerData) StartCameraDrag(point math.Vec2) {
	p.CameraDrag = &point
}

func (p *PlayerData) UpdateCameraDrag(point math.Vec2) (math.Vec2, bool) {
	if p.CameraDrag == nil {
		p.StartCameraDrag(point)
		return math.Vec2{}, false
	}

	delta := point.Sub(*p.CameraDrag).MulScalar(2)
	p.CameraDrag = &point
	if delta.IsZero() {
		return math.Vec2{}, false
	}

	return delta, true
}

func (p *PlayerData) ClearCameraDrag() {
	p.CameraDrag = nil
}

func (p *PlayerData) ScreenToWorld(point math.Vec2) math.Vec2 {
	if p.Camera == nil {
		return point
	}

	x, y := p.Camera.GetWorldCoords(point.X, point.Y)
	return math.NewVec2(x, y)
}

func (p *PlayerData) DragDistance() float64 {
	if p.DragStart == nil || p.DragEnd == nil {
		return 0
	}

	return p.DragStart.Distance(*p.DragEnd)
}
func (p *PlayerData) DragRect() image.Rectangle {
	if p.DragStart == nil || p.DragEnd == nil {
		return image.Rectangle{}
	}

	return util.ToRect(*p.DragStart, *p.DragEnd)
}

var Player = donburi.NewComponentType[PlayerData]()

func GetPlayerEntity(world donburi.World) donburi.Entity {
	entry := Player.MustFirst(world)
	return entry.Entity()
}
func GetPlayer(world donburi.World) *PlayerData {
	entry := Player.MustFirst(world)
	return Player.Get(entry)
}
