package components

import (
	"testing"

	dmath "github.com/yohamta/donburi/features/math"
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
			got := initialTargetLoopIndex(test.targets, test.loop, test.current, DEFAULT_STOP_DISTANCE)
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
