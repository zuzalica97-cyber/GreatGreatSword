package game

import "github.com/hajimehoshi/ebiten/v2"

type Entity interface {
	Update(world WorldView) bool
	Draw(screen *ebiten.Image)
	Tag() string
}
