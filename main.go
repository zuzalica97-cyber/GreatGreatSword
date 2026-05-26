package main

import (
	"embed"
	"fmt"
	"great-sword/game/common"
	"great-sword/game/player"
	swords "great-sword/game/player/Swords"
	game "great-sword/game/world"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	mplusFaceSource *text.GoTextFaceSource
	gameAssets      embed.FS
)

type Game struct {
	world        *game.World
	gameOver     bool
	lowerTexture *ebiten.Image
}

func NewGame() *Game {
	var err error

	g := &Game{}

	g.lowerTexture, _, err = ebitenutil.NewImageFromFile("assets/pod.png") //НЕ РАБОТАЕТ СДЕСЬ
	if err != nil {
		log.Fatal("failed to load lower texture", err)
	}

	return g
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

	if g.lowerTexture == nil {
		fmt.Println("olol")
	}

	if g.lowerTexture != nil {
		opBg := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.lowerTexture, opBg)
	} else {
		screen.Fill(color.RGBA{30, 30, 30, 255})
	}

	for _, entityH := range g.world.SearchEntities("playerHead") {
		p := entityH.(*player.PlayerHead)

		player.DrawRotatedRect(screen, p.PositionHead.Px, p.PositionHead.Py, common.PlayerHeadSize, common.PlayerHeadSize, p.Angle, color.RGBA{0, 0, 0, 255})
		player.DrawRotatedRect(screen, p.PositionHead.Px, p.PositionHead.Py, common.PlayerHeadSize, common.PlayerHeadSize, p.Angle, color.RGBA{255, 185, 185, 255})
	}
	for _, entity := range g.world.SearchEntities("blueSword") {
		s := entity.(*swords.BlueSword)

		player.DrawRotatedRect(screen, s.Position.Px, s.Position.Py, common.SwordAttachmentWidth, common.SwordAttachmentHeight, s.Angle, color.RGBA{0, 100, 200, 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth, common.ScreenHeight
}

func main() {
	ebiten.SetWindowSize(common.ScreenWidth, common.ScreenHeight)
	ebiten.SetWindowTitle("GratGreatSword_v0.1")

	NewGame()

	world := game.NewWorld()

	world.AddEntity(
		player.NewPlayerLeg(),
	)
	world.AddEntity(
		player.NewPlayerHead(),
	)
	world.AddEntity(
		swords.NewBlueSword(world),
	)

	g := &Game{
		world: world,
	}
	fmt.Println("g")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
