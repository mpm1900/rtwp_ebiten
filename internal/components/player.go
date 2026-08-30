package components

import (
	"image"
	"rtwp_ebitengine/internal/util"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

type Ability struct {
	AbilityID uuid.UUID
	Key       ebiten.Key
	Name      string
	Handle    func(*ecs.ECS)
}

type PlayerData struct {
	Ability   *Ability
	DragStart *math.Vec2
	DragEnd   *math.Vec2
}

func NewPlayerData(ability *Ability) PlayerData {
	return PlayerData{
		Ability:   ability,
		DragStart: nil,
		DragEnd:   nil,
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
