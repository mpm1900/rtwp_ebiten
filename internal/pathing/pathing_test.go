package pathing

import (
	"rtwp_ebitengine/internal/components"
	"testing"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func createObstacle(world donburi.World, pos, size dmath.Vec2) donburi.Entity {
	entity := world.Create(components.Collision, transform.Transform)
	entry := world.Entry(entity)
	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: pos,
		LocalScale:    size,
	})
	return entity
}

func TestFindPath_Direct(t *testing.T) {
	world := donburi.NewWorld()
	start := dmath.NewVec2(200, 200)
	goal := dmath.NewVec2(300, 200)

	path, ok := FindPath(world, start, goal)
	if !ok {
		t.Fatalf("expected path to be found, got ok=false")
	}
	if len(path) == 0 {
		t.Fatalf("expected non-empty path")
	}
	lastPoint := path[len(path)-1]
	if lastPoint.Distance(goal) > 0.001 {
		t.Errorf("expected final path point to be %v, got %v", goal, lastPoint)
	}
}

func TestFindPath_Diagonal_Subdivided(t *testing.T) {
	world := donburi.NewWorld()
	start := dmath.NewVec2(200, 200)
	goal := dmath.NewVec2(500, 500)

	path, ok := FindPath(world, start, goal)
	if !ok {
		t.Fatalf("expected path to be found, got ok=false")
	}
	// Distance is ~424px. With MaxPointDistance=32, expect ~14 points along diagonal
	if len(path) < 10 {
		t.Errorf("expected dense waypoints along diagonal, got %d points: %v", len(path), path)
	}
	lastPoint := path[len(path)-1]
	if lastPoint.Distance(goal) > 0.001 {
		t.Errorf("expected path end to be %v, got %v", goal, lastPoint)
	}
}

func TestFindPath_SamePosition(t *testing.T) {
	world := donburi.NewWorld()
	start := dmath.NewVec2(200, 200)
	goal := dmath.NewVec2(200, 200)

	path, ok := FindPath(world, start, goal)
	if !ok {
		t.Fatalf("expected path to be found, got ok=false")
	}
	if len(path) != 1 || path[0] != goal {
		t.Errorf("expected path to be [goal], got %v", path)
	}
}

func TestFindPath_AroundObstacle(t *testing.T) {
	world := donburi.NewWorld()
	createObstacle(world, dmath.NewVec2(240, 180), dmath.NewVec2(48, 48))

	start := dmath.NewVec2(200, 200)
	goal := dmath.NewVec2(350, 200)

	path, ok := FindPath(world, start, goal)
	if !ok {
		t.Fatalf("expected path to navigate around obstacle, got ok=false")
	}
	if len(path) == 0 {
		t.Fatalf("expected non-empty path")
	}
	lastPoint := path[len(path)-1]
	if lastPoint.Distance(goal) > 0.001 {
		t.Errorf("expected final path point to be %v, got %v", goal, lastPoint)
	}
}

func TestFindPath_StartInsideOwnEntity(t *testing.T) {
	world := donburi.NewWorld()
	// Moving entity at (200, 200) with size (24, 24)
	createObstacle(world, dmath.NewVec2(188, 188), dmath.NewVec2(24, 24))

	start := dmath.NewVec2(200, 200)
	goal := dmath.NewVec2(400, 200)

	path, ok := FindPath(world, start, goal)
	if !ok {
		t.Fatalf("expected path to be found when starting inside own collision entity, got ok=false")
	}
	lastPoint := path[len(path)-1]
	if lastPoint.Distance(goal) > 0.001 {
		t.Errorf("expected final path point to be %v, got %v", goal, lastPoint)
	}
}

func BenchmarkFindPath_DirectLOS(b *testing.B) {
	world := donburi.NewWorld()
	start := dmath.NewVec2(200, 200)
	goal := dmath.NewVec2(600, 600)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindPath(world, start, goal)
	}
}

func BenchmarkFindPath_AroundObstacles(b *testing.B) {
	world := donburi.NewWorld()
	createObstacle(world, dmath.NewVec2(240, 180), dmath.NewVec2(48, 48))
	createObstacle(world, dmath.NewVec2(500, 500), dmath.NewVec2(100, 100))

	start := dmath.NewVec2(200, 200)
	goal := dmath.NewVec2(600, 600)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindPath(world, start, goal)
	}
}
