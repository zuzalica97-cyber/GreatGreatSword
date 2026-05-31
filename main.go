package main

import (
	"embed"
	"fmt"
	enemy "great-sword/game/Enemy"
	"great-sword/game/common"
	"great-sword/game/player"
	swords "great-sword/game/player/Swords"
	game "great-sword/game/world"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	mplusFaceSource *text.GoTextFaceSource
	gameAssets      embed.FS
)

type Game struct {
	world        *game.World
	gameOver     bool
	lowerTexture *ebiten.Image
	playerTex    *ebiten.Image
}

func NewGame(world *game.World) *Game {

	var err error

	g := &Game{
		world: world,
	}

	g.lowerTexture, _, err = ebitenutil.NewImageFromFile("assets/poll.png") //НЕ РАБОТАЕТ СДЕСЬ
	if err != nil {
		log.Fatal("failed to load lower texture", err)
	}
	g.playerTex, _, err = ebitenutil.NewImageFromFile("assets/goll.png") //НЕ РАБОТАЕТ СДЕСЬ
	if err != nil {
		log.Fatal("failed to load lower texture", err)
	}

	return g
}

func (g *Game) ResetGame() {
	common.Score = 0

	for _, entity := range g.world.SearchEntities("playerLeg") {
		leg := entity.(*player.PlayerLeg)

		leg.Position.Px = common.ScreenWidth/2 - common.PlayerSize/2
		leg.Position.Py = common.ScreenHeight/2 - common.PlayerSize/2
		leg.Speed.Vx = 0
		leg.Speed.Vy = 0

		for _, entity := range g.world.SearchEntities("playerHead") {
			head := entity.(*player.PlayerHead)

			head.Angle = 0
			head.AngularVelocity = 0

			for _, entity := range g.world.SearchEntities("blueSword") {
				sword := entity.(*swords.BlueSword)

				sword.UpdateAttachmentTarget(g.world)
				sword.Position.Px = sword.TargetX
				sword.Position.Py = sword.TargetY
				sword.Angle = sword.TargetAngle

				for _, entity := range g.world.SearchEntities("pathetic") {
					path := entity.(*enemy.Pathetic)

					path.Paths = make([]enemy.OnePath, 0)
				}
			}
		}
	}

}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	if common.PlayerHelth <= 0 {
		g.ResetGame()
		common.PlayerHelth = common.MaxPlayerHelth
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

	if g.lowerTexture != nil {
		op := &ebiten.DrawImageOptions{}

		texWidth := float64(g.lowerTexture.Bounds().Dx())
		texHeidth := float64(g.lowerTexture.Bounds().Dy())

		scaleX := common.ScreenWidth / texWidth
		scaleY := common.ScreenHeight / texHeidth

		op.GeoM.Scale(scaleX, scaleY)

		screen.DrawImage(g.lowerTexture, op)
	} else {
		screen.Fill(color.RGBA{30, 30, 30, 255})
	}

	for _, entityH := range g.world.SearchEntities("playerHead") {
		p := entityH.(*player.PlayerHead)

		centeX := p.PositionHead.Px + common.PlayerHeadSize/2
		centeY := p.PositionHead.Py + common.PlayerHeadSize/2

		if g.playerTex != nil {

			texW := float64(g.playerTex.Bounds().Dx())
			texH := float64(g.playerTex.Bounds().Dy())

			scaleX := common.PlayerHeadSize / texW
			sclalY := common.PlayerHeadSize / texH

			offsetX := -20.0
			offsetY := -20.0

			opHat := &ebiten.DrawImageOptions{}
			opHat.GeoM.Translate(-texW/2, -texH/2)
			opHat.GeoM.Scale(scaleX, sclalY)
			opHat.GeoM.Rotate(p.Angle * math.Pi / 180)
			opHat.GeoM.Translate(centeX+offsetX, centeY+offsetY) //ДОДЕЛАТЬ
			screen.DrawImage(g.playerTex, opHat)
		}
	}
	for _, entity := range g.world.SearchEntities("blueSword") {
		s := entity.(*swords.BlueSword)

		player.DrawRotatedRect(screen, s.Position.Px, s.Position.Py, common.SwordAttachmentWidth, common.SwordAttachmentHeight, s.Angle, color.RGBA{0, 100, 200, 255})
	}
	for _, entity := range g.world.SearchEntities("pathetic") {
		p := entity.(*enemy.Pathetic)

		for _, enemy := range p.Paths {
			vector.FillRect(
				screen,
				float32(enemy.PX),
				float32(enemy.PY),
				common.PatheticSize,
				common.PatheticSize,
				color.RGBA{150, 150, 150, 255},
				true,
			)

		}

	}

	for _, entity := range g.world.SearchEntities("hater") {
		h := entity.(*enemy.Haters)

		for _, enemy := range h.HatersMass {
			vector.FillRect(
				screen,
				float32(enemy.HX),
				float32(enemy.HY),
				common.HaterSize,
				common.HaterSize,
				color.RGBA{180, 150, 150, 255},
				true,
			)

		}

	}

	scoreStr := fmt.Sprint("SCORE: ", common.Score)
	textIng := ebiten.NewImage(400, 80)
	ebitenutil.DebugPrintAt(textIng, scoreStr, 10, 10)

	scoreOp := &ebiten.DrawImageOptions{}
	scoreOp.GeoM.Scale(3.0, 3.0)
	scoreOp.GeoM.Translate(common.ScreenWidth-250, 10)
	screen.DrawImage(textIng, scoreOp)

	helthStr := fmt.Sprint("HELTH: ", common.PlayerHelth)
	helthIng := ebiten.NewImage(400, 80)
	ebitenutil.DebugPrintAt(helthIng, helthStr, 10, 10)

	helthOp := &ebiten.DrawImageOptions{}
	helthOp.GeoM.Scale(3.0, 3.0)
	helthOp.GeoM.Translate(common.ScreenWidth-250, 70)
	screen.DrawImage(helthIng, helthOp)
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
	world.AddEntity(
		swords.NewBlueSword(world),
	)
	world.AddEntity(
		enemy.NewPathetic(),
	)
	world.AddEntity(
		enemy.NewHaters(),
	)

	g := NewGame(world)

	fmt.Println("g")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
