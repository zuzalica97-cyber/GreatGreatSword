package main

import (
	"embed"
	"fmt"
	enemy "great-sword/game/Enemy"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
	"great-sword/game/player"
	swords "great-sword/game/player/Swords"
	game "great-sword/game/world"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	mplusFaceSource *text.GoTextFaceSource
	gameAssets      embed.FS
)

type Game struct {
	world                        *game.World
	gameOver                     bool
	lowerTexture                 *ebiten.Image
	MainKamera                   kamera.Camera
	cameraSmoothX, cameraSmoothY float64
	CollisionManager             *hitboxes.CollisionManager
	frameCounter                 float64 // ВРЕМЕННО
}

func NewGame(world *game.World, manager *hitboxes.CollisionManager) *Game {
	var g *Game
	var err error
	g = &Game{
		world:            world,
		MainKamera:       *kamera.NewCamera(0, 0, common.ScreenWidth, common.ScreenHeight),
		CollisionManager: manager,
		frameCounter:     1.0,
	}

	g.lowerTexture, _, err = ebitenutil.NewImageFromFile("assets/poll.png") //НЕ РАБОТАЕТ СДЕСЬ
	if err != nil {
		log.Fatal("failed to load lower texture", err)
	}

	return g
}

func (g *Game) ResetGame() {
	common.Score = 0

	for _, entity := range g.world.SearchEntities("playerLeg") {
		leg := entity.(*player.PlayerLeg)

		leg.Position.Px = common.RoomWidth/2 - common.PlayerSize/2
		leg.Position.Py = common.RoomHeight/2 - common.PlayerSize/2
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

					for i, _ := range path.Paths {
						g.CollisionManager.RemoveObject(path.Paths[i]) //Удаляем из обектов коллизии
					}

					path.Paths = make([]*enemy.OnePath, 0)

					for _, entity := range g.world.SearchEntities("hater") {
						hater := entity.(*enemy.Haters)

						for i, _ := range hater.HatersMass {
							g.CollisionManager.RemoveObject(hater.HatersMass[i]) //Удаляем из обектов коллизии
						}

						hater.HatersMass = make([]*enemy.Haits, 0)

					}
				}
			}
		}
	}

}

func (g *Game) Update() error { // ДЗ нужно сделать так чтобы игра не ломалась при большём количестве врогов. далее добаваить систему васа и икоолизию для стен и игрока. переделать или удолить некоторые функции с коолизией кторые были раньше. добавить противника с большим весом и оружием.
	if g.gameOver {
		return nil
	}

	if g.frameCounter > 0 {
		g.frameCounter -= 1.0 / 60.0
		if g.frameCounter <= 0 {
			g.frameCounter = 1.0
			g.CollisionManager.Cleanup()
		}
	}

	if common.PlayerHelth <= 0 {
		g.ResetGame()
		common.PlayerHelth = common.MaxPlayerHelth
	}

	for _, entity := range g.world.Entities() {
		if entity.Update(g.world, g.CollisionManager) {
			g.gameOver = true
			return nil
		}
	}

	g.CollisionManager.Update()

	for _, Pleg := range g.world.SearchEntities("playerLeg") {
		p := Pleg.(*player.PlayerLeg)

		// Стало (смотрит в ЦЕНТР):
		g.MainKamera.LookAt(
			p.Position.Px+common.PlayerSize/2,
			p.Position.Py+common.PlayerSize/2,
		)
		if ebiten.IsKeyPressed(ebiten.KeyL) {
			log.Println(p.Position.Px+common.PlayerSize/2,
				p.Position.Py+common.PlayerSize/2,
			)
		}
	}

	// РУЧНОЕ СГЛАЖИВАНИЕ (0.05 = плавность)
	g.cameraSmoothX += (g.MainKamera.X - g.cameraSmoothX) * 0.1
	g.cameraSmoothY += (g.MainKamera.Y - g.cameraSmoothY) * 0.1

	// Применяем сглаженные координаты
	g.MainKamera.X = g.cameraSmoothX
	g.MainKamera.Y = g.cameraSmoothY

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{40, 40, 60, 255})

	if g.lowerTexture != nil {
		op := &ebiten.DrawImageOptions{}

		texWidth := float64(g.lowerTexture.Bounds().Dx())
		texHeight := float64(g.lowerTexture.Bounds().Dy())

		scaleX := common.RoomWidth / texWidth
		scaleY := common.RoomHeight / texHeight

		op.GeoM.Scale(scaleX, scaleY)
		g.MainKamera.Draw(g.lowerTexture, op, screen)
	} else {
		screen.Fill(color.RGBA{30, 30, 30, 255})
	}

	for _, entity := range g.world.DrawlerEntities() {
		entity.Draw(screen, &g.MainKamera)
	}

	// ===== UI (НЕ ПРИМЕНЯЕМ КАМЕРУ) =====

	// Счёт
	scoreStr := fmt.Sprint("SCORE: ", common.Score)
	textIng := ebiten.NewImage(400, 80)
	ebitenutil.DebugPrintAt(textIng, scoreStr, 10, 10)

	scoreOp := &ebiten.DrawImageOptions{}
	scoreOp.GeoM.Scale(3.0, 3.0)
	scoreOp.GeoM.Translate(common.ScreenWidth-250, 10)
	screen.DrawImage(textIng, scoreOp)

	// Здоровье
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

	manager := hitboxes.NewCollisionManager()

	world.AddEntity(
		player.NewPlayerLeg(manager),
	)
	world.AddEntity(
		player.NewPlayerHead(),
	)
	world.AddEntity(
		swords.NewBlueSword(world, manager),
	)
	world.AddEntity(
		enemy.NewPathetic(),
	)
	world.AddEntity(
		enemy.NawHaters(),
	)
	world.AddEntity(
		enemy.NewBerserk(),
	)

	g := NewGame(world, manager)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
