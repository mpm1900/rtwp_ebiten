package components

import (
	"image"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
)

const collisionCellSize = 64

type collisionIndex struct {
	cells map[[2]int][]donburi.Entity
}

var collisionSpatial collisionIndex

func RebuildCollisionIndex(world donburi.World) {
	if collisionSpatial.cells == nil {
		collisionSpatial.cells = map[[2]int][]donburi.Entity{}
	}
	clear(collisionSpatial.cells)

	for entry := range CollisionQuery.Iter(world) {
		bounds, ok := Rect(entry)
		if !ok {
			continue
		}

		collisionSpatial.add(entry.Entity(), bounds)
	}
}

func (idx *collisionIndex) add(entity donburi.Entity, bounds image.Rectangle) {
	min_x := cellCoord(bounds.Min.X)
	min_y := cellCoord(bounds.Min.Y)
	max_x := cellCoord(bounds.Max.X - 1)
	max_y := cellCoord(bounds.Max.Y - 1)

	for x := min_x; x <= max_x; x++ {
		for y := min_y; y <= max_y; y++ {
			key := [2]int{x, y}
			idx.cells[key] = append(idx.cells[key], entity)
		}
	}
}

func (idx *collisionIndex) eachCandidate(
	world donburi.World,
	bounds image.Rectangle,
	skip donburi.Entity,
	yield func(*donburi.Entry) bool,
) {
	min_x := cellCoord(bounds.Min.X)
	min_y := cellCoord(bounds.Min.Y)
	max_x := cellCoord(bounds.Max.X - 1)
	max_y := cellCoord(bounds.Max.Y - 1)

	for x := min_x; x <= max_x; x++ {
		for y := min_y; y <= max_y; y++ {
			for _, entity := range idx.cells[[2]int{x, y}] {
				if entity == skip || !world.Valid(entity) {
					continue
				}

				other := world.Entry(entity)
				if !other.HasComponent(Collision) || !other.HasComponent(transform.Transform) {
					continue
				}

				if !yield(other) {
					return
				}
			}
		}
	}
}

func cellCoord(value int) int {
	return value / collisionCellSize
}
