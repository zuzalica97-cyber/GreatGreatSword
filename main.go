package main

import (
	"fmt"
	"great-sword/game/common"
	"great-sword/game/player"
	game "great-sword/game/world"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

type Game struct {
	world    *game.World
	gameOver bool
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	for _, entity := range g.world.Entities() {
		if entity.Update(g.world) {
			g.gameOver = true
			return nil
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, entity := range g.world.SearchEntities("playerLeg") {
		p := entity.(*player.PlayerLeg)

		screen.Fill(color.RGBA{30, 30, 30, 255})

		speed := math.Sqrt(p.Speed.Vx*p.Speed.Vx + p.Speed.Vy*p.Speed.Vy)
		intensity := uint8(math.Min(255, speed/2))

		vector.FillRect(
			screen,
			float32(p.Position.Px),
			float32(p.Position.Py),
			common.PlayerSize,
			common.PlayerSize,
			color.RGBA{0, 255, 0, intensity},
			true,
		)

		vector.FillRect(
			screen,
			float32(p.Position.Px),
			float32(p.Position.Py),
			common.PlayerSize,
			common.PlayerSize,
			color.RGBA{255, 255, 255, 255},
			true,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth, common.ScreenHeight
}

func main() {
	ebiten.SetWindowSize(common.ScreenWidth, common.ScreenHeight)
	ebiten.SetWindowTitle("GratGreatSword_v0.1")

	world := game.NewWorld()

	world.AddEntity(
		player.NewPlayerLeg(),
	)

	g := &Game{
		world: world,
	}
	fmt.Println("g")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
