package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
)

func HandlePause(ecs *ecs.ECS) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if ecs.IsPaused() {
			fmt.Println("resume")
			ecs.Resume()
		} else {
			fmt.Println("pause")
			ecs.Pause()
		}
	}
}
