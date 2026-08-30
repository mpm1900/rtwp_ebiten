package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"

	camera "github.com/melonfunction/ebiten-camera"
)

type Ability struct {
	AbilityID uuid.UUID
	Key       ebiten.Key
	Name      string
	Handle    func(*ecs.ECS)
}

type PlayerData struct {
	Ability    *Ability
	Camera     *camera.Camera
	CameraDrag *math.Vec2
	DragStart  *math.Vec2
	DragEnd    *math.Vec2
}

func NewPlayerData() PlayerData {
	return PlayerData{
		Ability:    nil,
		Camera:     camera.NewCamera(640, 480, 320, 240, 0, 1),
		CameraDrag: nil,
		DragStart:  nil,
		DragEnd:    nil,
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

	delta := point.Sub(*p.CameraDrag)
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

func GetPlayer(world donburi.World) *PlayerData {
	entry, ok := Player.First(world)
	if !ok {
		return nil
	}

	return Player.Get(entry)
}
