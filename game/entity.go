package game

import (
	"great-sword/game/hitboxes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

type Entity interface {
	Update(world WorldView, manager *hitboxes.CollisionManager) bool
	Draw(screen *ebiten.Image, camera *kamera.Camera)
	Tag() string
}
