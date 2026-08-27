package util

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/component"
)

type componentKey struct {
	Entity    donburi.Entity
	Component component.ComponentTypeId
}

type Frame struct {
	originals map[componentKey]func(donburi.World)
}

type frameCloneable[T any] interface {
	Clone() T
}

func NewFrame() *Frame {
	return &Frame{
		originals: make(map[componentKey]func(donburi.World)),
	}
}

func (f *Frame) Modify[T any](
	entry *donburi.Entry,
	componentType *donburi.ComponentType[T],
	updater func(*T),
) {
	if f.originals == nil {
		f.originals = make(map[componentKey]func(donburi.World))
	}

	key := componentKey{
		Entity:    entry.Entity(),
		Component: componentType.Id(),
	}

	value := componentType.Get(entry)
	if _, ok := f.originals[key]; !ok {
		original := cloneFrameValue(*value)
		f.originals[key] = func(world donburi.World) {
			if !world.Valid(key.Entity) {
				return
			}

			restoredEntry := world.Entry(key.Entity)
			if !restoredEntry.HasComponent(componentType) {
				return
			}

			componentType.SetValue(restoredEntry, original)
		}
	}

	updater(value)
}

func cloneFrameValue[T any](value T) T {
	if cloneable, ok := any(value).(frameCloneable[T]); ok {
		return cloneable.Clone()
	}
	return value
}

func (f *Frame) Restore(world donburi.World) {
	for _, restore := range f.originals {
		restore(world)
	}
	clear(f.originals)
}
