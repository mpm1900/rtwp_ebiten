package components

import (
	"testing"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func TestInitialTargetLoopIndex(t *testing.T) {
	current := dmath.NewVec2(100, 100)
	nearby := dmath.NewVec2(100.5, 100.5)
	far := dmath.NewVec2(200, 200)
	origin := dmath.NewVec2(100, 100)

	tests := []struct {
		name     string
		targets  []dmath.Vec2
		loop     bool
		current  dmath.Vec2
		expected int
	}{
		{
			name:     "non loop starts at first target",
			targets:  []dmath.Vec2{far},
			loop:     false,
			current:  current,
			expected: 0,
		},
		{
			name:     "loop skips nearby first waypoint",
			targets:  []dmath.Vec2{nearby, far, origin},
			loop:     true,
			current:  current,
			expected: 1,
		},
		{
			name:     "loop keeps distant first waypoint",
			targets:  []dmath.Vec2{far, origin},
			loop:     true,
			current:  current,
			expected: 0,
		},
		{
			name:     "loop single target starts at origin",
			targets:  []dmath.Vec2{origin},
			loop:     true,
			current:  current,
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := initialLoopIndex(test.targets, test.loop, test.current, DEFAULT_STOP_DISTANCE)
			if got != test.expected {
				t.Fatalf("expected index %d, got %d", test.expected, got)
			}
		})
	}
}

func TestAppendLoopOrigin(t *testing.T) {
	origin := dmath.NewVec2(100, 100)
	goal := dmath.NewVec2(300, 300)

	path := AppendLoopOrigin([]dmath.Vec2{goal}, origin)
	if len(path) != 2 {
		t.Fatalf("expected origin appended, got %v", path)
	}
	if path[1] != origin {
		t.Fatalf("expected origin %v, got %v", origin, path[1])
	}

	path = AppendLoopOrigin([]dmath.Vec2{goal, origin}, origin)
	if len(path) != 2 {
		t.Fatalf("expected duplicate origin skipped, got %v", path)
	}
}

func TestMovementLastIgnoresLoopOrigin(t *testing.T) {
	origin := dmath.NewVec2(100, 100)
	goal := dmath.NewVec2(300, 300)
	movement := MovementData{
		Loop:          true,
		LoopOrigin:    origin,
		HasLoopOrigin: true,
		Path:          []dmath.Vec2{goal, origin},
	}

	last, ok := movement.Last()
	if !ok || last != goal {
		t.Fatalf("expected last queued point %v, got %v, %t", goal, last, ok)
	}
}

func TestPushMovementListExtendsLoopBeforeOrigin(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(transform.Transform)
	entry := world.Entry(entity)
	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: dmath.NewVec2(90, 90),
		LocalScale:    dmath.NewVec2(20, 20),
	})

	origin := Center(entry)
	first_goal := dmath.NewVec2(300, 300)
	second_goal := dmath.NewVec2(500, 300)

	WithMovement(entry, NewPathMovement(entry, []dmath.Vec2{first_goal}, true))
	PushMovementList(entry, []dmath.Vec2{second_goal}, DEFAULT_STOP_DISTANCE, true)

	movement := Movement.Get(entry)
	expected_path := []dmath.Vec2{first_goal, second_goal, origin}
	if len(movement.Path) != len(expected_path) {
		t.Fatalf("expected path %v, got %v", expected_path, movement.Path)
	}

	for i, expected := range expected_path {
		if movement.Path[i] != expected {
			t.Fatalf("expected path %v, got %v", expected_path, movement.Path)
		}
	}
}
