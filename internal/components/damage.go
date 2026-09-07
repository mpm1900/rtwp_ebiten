package components

import (
	"fmt"
	"math/rand/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
)

var Damage = donburi.NewComponentType[float64]()
var DamageQuery = donburi.NewQuery(filter.Contains(Damage, Stats))

func GetHealth(entry *donburi.Entry) (health float64, damage float64) {
	stats := Stats.Get(entry)
	damage = 0.0
	if entry.HasComponent(Damage) {
		damage = *Damage.Get(entry)
	}

	return stats.Stats[StatHealth], damage
}

func WithDamage(entry *donburi.Entry, damage float64) {
	if !entry.HasComponent(Damage) {
		entry.AddComponent(Damage)
	}
	Damage.SetValue(entry, damage)
}

func DamageAt(world donburi.World, source donburi.Entity, point math.Vec2, action *Action) (*donburi.Entry, bool) {
	entry, ok := FirstActorAtPoint(world, point)
	if !ok {
		return nil, false
	}

	source_stats := Stats.Get(world.Entry(source))
	target_stats := Stats.Get(entry)
	result := GetDamageResult(action, 20, *source_stats, *target_stats)
	if !entry.HasComponent(Damage) {
		WithDamage(entry, result.Amount)
		return nil, false
	}

	damage := Damage.Get(entry)
	Damage.SetValue(entry, *damage+result.Amount)
	return entry, true
}

type AccuracyResult struct {
	BaseAccuracy float64
	Accuracy     float64
	Roll         float64
	Success      bool
}

type DamageResult struct {
	AccuracyResult AccuracyResult
	Amount         float64
}

func GetAccuracyResult(action *Action, source StatsData, target StatsData) AccuracyResult {
	stage := source.Stages[StatAccuracy] - target.Stages[StatEvasion]
	accuracy := action.Accuracy * getStageMult(stage, 3)
	roll := rand.Float64() * 100
	return AccuracyResult{
		BaseAccuracy: action.Accuracy,
		Accuracy:     accuracy,
		Roll:         roll,
		Success:      accuracy > roll,
	}
}

func GetDamageResult(action *Action, level int, source StatsData, target StatsData) DamageResult {
	accuracy_result := GetAccuracyResult(action, source, target)
	ratio := source.Stats[StatMelee] / target.Stats[StatDefense]
	level_mod := float64(level*2)/5 + 2
	amount := (action.Power*ratio*level_mod)/50 + 2
	if !accuracy_result.Success {
		fmt.Println("MISS")
		amount = 0
	} else {
		fmt.Println("HIT")
	}

	return DamageResult{
		AccuracyResult: accuracy_result,
		Amount:         amount,
	}
}
