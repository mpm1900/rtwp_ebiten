package ecs

import "github.com/yohamta/donburi"

var Delay = donburi.NewComponentType[int]()

func WithDelay(entry *donburi.Entry, delay int) {
	entry.AddComponent(Delay)
	Delay.SetValue(entry, delay)
}

func DecrementDelays(world donburi.World) {
	for entry := range Delay.Iter(world) {
		delay := Delay.Get(entry)
		if *delay > 0 {
			*delay--
		}
	}
}

func RemoveCompletedDelays(world donburi.World) {
	for entry := range Delay.Iter(world) {
		delay := *Delay.Get(entry)
		if delay == 0 {
			entry.RemoveComponent(Delay)
		}
	}
}
