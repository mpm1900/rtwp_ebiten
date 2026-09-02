package components

import (
	"github.com/yohamta/donburi"
)

var Delay = donburi.NewComponentType[int]()

func WithDelay(entry *donburi.Entry, delay int) {
	if !entry.HasComponent(Delay) {
		entry.AddComponent(Delay)
	}

	Delay.SetValue(entry, delay)
}
