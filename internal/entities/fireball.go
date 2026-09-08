package entities

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

const fireballSpeed = 80.0

func CreateFireball(ecs *ecs.ECS, source donburi.Entity, target donburi.Entity, position math.Vec2) donburi.Entity {
	entity := ecs.Create(EffectLayer, components.Image, components.Movement)
	entry := ecs.World.Entry(entity)
	components.WithImage(entry, assets.FireImage, position)
	components.WithMovement(entry, components.NewHomingMovement(target, fireballSpeed))
	components.WithCollision(entry, components.CollisionData{
		OnCollide: func(w donburi.World, e1 *donburi.Entry, e2 donburi.Entity) {
			if !w.Valid(e2) {
				return
			}

			other := w.Entry(e2)
			if !other.HasComponent(components.Actor) {
				return
			}

			components.DamageTo(w, source, e2, 100, 80)
			events.Destroy.Publish(w, e1.Entity())
		},
	})
	return entity
}
