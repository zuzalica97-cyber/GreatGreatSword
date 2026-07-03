package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

type Entity interface {
	Update(world WorldView) bool
	Draw(screen *ebiten.Image, camera *kamera.Camera)
	Tag() string
}
