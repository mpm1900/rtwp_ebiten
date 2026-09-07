package modifiers

import (
	"maps"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

type ResolveBaseStatsBehavior struct{}

func (b ResolveBaseStatsBehavior) Active(mod *components.ModifierData, world donburi.World, modifier *donburi.Entry) bool {
	return true
}
func (b ResolveBaseStatsBehavior) Apply(mod *components.ModifierData, world donburi.World, frame *util.Frame, modifier *donburi.Entry) {
	components.EachDependent(world, modifier, func(entry *donburi.Entry) {
		if entry.HasComponent(components.Stats) {
			frame.Modify(entry, components.Stats, func(stats *components.StatsData) {
				stats.Stats = maps.Clone(stats.Base)
				stats.ResolveBaseStats()
				stats.MapStages()
			})
		}
	})
}
func (b ResolveBaseStatsBehavior) Spawn(mod *components.ModifierData, ecs *ecs.ECS, position math.Vec2) donburi.Entity {
	entity := entities.CreateEffect(ecs, mod)
	entry := ecs.World.Entry(entity)
	entry.AddComponent(components.TargetsWhere)
	components.TargetsWhere.SetValue(entry, func(e donburi.Entity) bool {
		return true
	})

	return entity
}

var SystemResolveBaseStats = &components.ModifierData{
	Priority: 1,
	Behavior: ResolveBaseStatsBehavior{},
}
