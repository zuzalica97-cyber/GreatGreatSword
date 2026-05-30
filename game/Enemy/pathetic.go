package enemy

import (
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/player"
	swords "great-sword/game/player/Swords"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Entity = (*Pathetic)(nil)

type Pathetic struct {
	Paths []OnePath
}

type OnePath struct {
	PX, PY float64
	Active bool
}

func NewPathetic() *Pathetic {
	p := &Pathetic{}
	return p
}

func (p *Pathetic) SpawnPathetic() {
	side := rand.IntN(4)

	var x, y float64

	switch side {
	case 0: //левый крвй
		x = -float64(common.PatheticSize)
		y = float64(rand.IntN(common.ScreenHeight))
	case 1: //левый крвй
		x = float64(common.ScreenWidth)
		y = float64(rand.IntN(common.ScreenHeight))
	case 2: //левый крвй
		x = float64(rand.IntN(common.ScreenWidth))
		y = -float64(common.PatheticSize)
	case 3: //левый крвй
		x = float64(rand.IntN(common.ScreenWidth))
		y = float64(common.ScreenHeight)
	}

	p.Paths = append(p.Paths, OnePath{
		PX:     x,
		PY:     y,
		Active: true,
	})
}

func (p *Pathetic) CheckCollisionWithAttachment(enemyX, enemyY, enemySize float64, attX, attY, attW, attH, attAngle float64) bool {
	cos := math.Cos(attAngle)
	sin := math.Sin(attAngle)

	enemyCenterX := enemyX + enemySize/2
	enemyCenterY := enemyY + enemySize/2

	localX := (enemyCenterX-attX)*cos + (enemyCenterY-attY)*sin
	localY := -(enemyCenterX-attX)*sin + (enemyCenterY-attY)*cos

	halfW := attW / 2
	halfH := attH / 2

	return math.Abs(localX) < halfW+enemySize/2 && math.Abs(localY) < halfH+enemySize/2
}

func (p *Pathetic) CheckCollisionWithPlayer(enemyX, enemyY, enemySize float64, world game.WorldView) bool {

	var playerX float64
	var playerY float64

	for _, entity := range world.SearchEntities("playerLeg") {
		leg := entity.(*player.PlayerLeg)

		playerX = leg.Position.Px
		playerY = leg.Position.Py

	}
	return enemyX < playerX+common.PlayerSize &&
		enemyX+enemySize > playerX &&
		enemyY < playerY+common.PlayerSize &&
		enemyY+enemySize > playerY
}

func (p *Pathetic) ResetGame(world game.WorldView) {
	common.Score = 0

	for _, entity := range world.SearchEntities("playerLeg") {
		leg := entity.(*player.PlayerLeg)

		leg.Position.Px = common.ScreenWidth/2 - common.PlayerSize/2
		leg.Position.Py = common.ScreenHeight/2 - common.PlayerSize/2
		leg.Speed.Vx = 0
		leg.Speed.Vy = 0

		for _, entity := range world.SearchEntities("playerHead") {
			head := entity.(*player.PlayerHead)

			head.Angle = 0
			head.AngularVelocity = 0
			p.Paths = make([]OnePath, 0)

			for _, entity := range world.SearchEntities("blueSword") {
				sword := entity.(*swords.BlueSword)

				sword.UpdateAttachmentTarget(world)
				sword.Position.Px = sword.TargetX
				sword.Position.Py = sword.TargetY
				sword.Angle = sword.TargetAngle

			}
		}
	}

}

func (p *Pathetic) UpatePathetics(dt float64, worldView game.WorldView) {
	for _, entity := range worldView.SearchEntities("playerLeg") {
		pLeg := entity.(*player.PlayerLeg)

		playerX := pLeg.Position.Px + common.PlayerSize/2
		playerY := pLeg.Position.Py + common.PlayerSize/2

		for _, entity := range worldView.SearchEntities("blueSword") {
			sword := entity.(*swords.BlueSword)

			attahmentX := sword.Position.Px
			attahmentY := sword.Position.Py
			attahmentW := common.SwordAttachmentWidth
			attahmentH := common.SwordAttachmentHeight
			attahmentAngle := sword.Angle * math.Pi / 180

			for i := 0; i < len(p.Paths); i++ {
				enemy := &p.Paths[i]
				if !enemy.Active {
					continue
				}

				dx := playerX - enemy.PX
				dy := playerY - enemy.PY
				distance := math.Sqrt(dx*dx + dy*dy)

				var speed float64

				if distance > common.PatheticDistanse {
					speed = common.PatheticBaseSpeed
				} else {
					t := 1.0 - distance/common.PatheticDistanse
					speed = common.PatheticBaseSpeed + t*(common.PatheticMaxSpeed-common.PatheticBaseSpeed)
				}
				if distance > 0.01 {
					enemy.PX += (dx / distance) * speed * dt
					enemy.PY += (dy / distance) * speed * dt
				}

				if p.CheckCollisionWithAttachment(enemy.PX, enemy.PY, common.PatheticSize, attahmentX,
					attahmentY, float64(attahmentW), float64(attahmentH), attahmentAngle) {
					p.Paths = append(p.Paths[:i], p.Paths[i+1:]...)
					i--
					common.Score++
					continue
				}
				if p.CheckCollisionWithPlayer(enemy.PX, enemy.PY, common.PatheticSize, worldView) {

					log.Println("Игрок погиб!")

					p.ResetGame(worldView) //ДОДЕЛАТЬ ФУНЦИЮ

				}
			}
		}
	}

}

func (p *Pathetic) Update(worldView game.WorldView) bool {
	dt := 1.0 / 60.0

	if len(p.Paths) < 5 {
		p.SpawnPathetic()
	}

	p.UpatePathetics(dt, worldView)

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		p.ResetGame(worldView)
	}
	return false
}

func (p *Pathetic) Draw(screen *ebiten.Image) {
}

func (p *Pathetic) Tag() string {
	return "pathetic"
}
