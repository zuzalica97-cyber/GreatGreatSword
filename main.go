package main

import (
	"fmt"
	"great-sword/game/common"
	"great-sword/game/player"
	game "great-sword/game/world"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

		centerX := p.Position.Px + common.PlayerSize/2
		centerY := p.Position.Py + common.PlayerSize/2

		for _, PlauerHe := range g.world.SearchEntities("playerHead") {
			head := PlauerHe.(*player.PlayerHead)

			player.DrawRotatedRect(screen, centerX, centerY, common.PlayerHeadSize, common.PlayerHeadSize, head.Angle, color.RGBA{200, 0, 0, 255})
			player.DrawRotatedRect(screen, centerX, centerY, common.PlayerHeadSize, common.PlayerHeadSize, head.Angle, color.RGBA{255, 180, 180, 255})

		}
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
	world.AddEntity(
		player.NewPlayerHead(),
	)

	g := &Game{
		world: world,
	}
	fmt.Println("g")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
