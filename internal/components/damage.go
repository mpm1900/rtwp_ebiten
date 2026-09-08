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

func IsAlive(entry *donburi.Entry) bool {
	if !entry.HasComponent(Stats) {
		return true
	}

	health, damage := GetHealth(entry)
	return health > damage
}

func WithDamage(entry *donburi.Entry, damage float64) {
	if !entry.HasComponent(Damage) {
		entry.AddComponent(Damage)
	}
	Damage.SetValue(entry, damage)
}

func DamageAt(world donburi.World, source donburi.Entity, point math.Vec2, accuracy, power float64) (*donburi.Entry, bool) {
	entry, ok := FirstActorAtPoint(world, point)
	if !ok {
		return nil, false
	}

	return DamageTo(world, source, entry.Entity(), accuracy, power)
}
func DamageTo(world donburi.World, source donburi.Entity, target donburi.Entity, accuracy, power float64) (*donburi.Entry, bool) {
	if !world.Valid(source) || !world.Valid(target) {
		return nil, false
	}

	source_entry := world.Entry(source)
	target_entry := world.Entry(target)
	if !source_entry.HasComponent(Stats) || !target_entry.HasComponent(Stats) {
		return nil, false
	}

	source_stats := Stats.Get(source_entry)
	target_stats := Stats.Get(target_entry)
	result := GetDamageResult(accuracy, power, 20, *source_stats, *target_stats)
	if !target_entry.HasComponent(Damage) {
		WithDamage(target_entry, result.Amount)
		return nil, false
	}

	damage := Damage.Get(target_entry)
	Damage.SetValue(target_entry, *damage+result.Amount)
	return target_entry, true
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

func GetAccuracyResult(base_accuracy float64, source StatsData, target StatsData) AccuracyResult {
	stage := source.Stages[StatAccuracy] - target.Stages[StatEvasion]
	accuracy := base_accuracy * getStageMult(stage, 3)
	roll := rand.Float64() * 100
	return AccuracyResult{
		BaseAccuracy: base_accuracy,
		Accuracy:     accuracy,
		Roll:         roll,
		Success:      accuracy > roll,
	}
}

func GetDamageResult(accuracy, power float64, level int, source StatsData, target StatsData) DamageResult {
	accuracy_result := GetAccuracyResult(accuracy, source, target)
	ratio := source.Stats[StatMelee] / target.Stats[StatDefense]
	level_mod := float64(level*2)/5 + 2
	amount := (power*ratio*level_mod)/50 + 2
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
