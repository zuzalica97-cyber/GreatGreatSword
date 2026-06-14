package enemy

import (
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/player"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Entity = (*Haters)(nil)

type Haters struct {
	HatersMass []Haits
	Bullets    []HaitersBullet
}

type Haits struct {
	HatersCooldownActive bool
	HatersCooldownTimer  float64
	BulletShotActive     bool
	BulletShotTimer      float64
	HX, HY               float64
	Active               bool
}
type HaitersBullet struct {
	BX, BY    float64
	VX, VY    float64
	BulActive bool
}

func NawHaters() *Haters {
	h := &Haters{}
	return h
}

func (h *Haters) SpawnHiters(x, y float64) {

	h.HatersMass = append(h.HatersMass, Haits{
		HX:     x,
		HY:     y,
		Active: true,
	})
}

func (h *Haters) CreateBullet(enemyX, enemyY, playerX, playerY float64, world game.WorldView) {

	dx := playerX - enemyX
	dy := playerY - enemyY
	diststance := math.Sqrt(dx*dx + dy*dy)

	if diststance > 0.01 {
		vx := (dx / diststance) * common.HaitersBolletSpeed
		vy := (dy / diststance) * common.HaitersBolletSpeed

		h.Bullets = append(h.Bullets, HaitersBullet{
			BX:        enemyX,
			BY:        enemyY,
			VX:        vx,
			VY:        vy,
			BulActive: true,
		})
	}
}

func (h *Haters) HaterShot(number int, enemyX, enemyY, playerX, playerY float64, world game.WorldView) {
	enemy := &h.HatersMass[number]

	enemy.BulletShotActive = true
	enemy.BulletShotTimer = common.HaitersBolletCooldown
	h.CreateBullet(enemyX, enemyY, playerX, playerY, world)
}

func (h *Haters) UpdateHaiters(dt float64, world game.WorldView) {

	playerX, playerY, attahmentX, attahmentY, attahmentW, attahmentH, attahmentAngle := WhereThePlayer(world)

	for i := 0; i < len(h.HatersMass); i++ {
		enemy := &h.HatersMass[i]
		if !enemy.Active {
			continue
		}

		if enemy.BulletShotActive {
			enemy.BulletShotTimer -= dt
			if enemy.BulletShotTimer <= 0 && !enemy.HatersCooldownActive {
				enemy.BulletShotActive = false
			}
		}

		if enemy.HatersCooldownActive {
			enemy.HatersCooldownTimer -= dt
			if enemy.HatersCooldownTimer <= 0 {
				enemy.HatersCooldownActive = false
			}
		}

		dx := playerX - enemy.HX
		dy := playerY - enemy.HY
		distance := math.Sqrt(dx*dx + dy*dy)

		var speed float64

		if distance > common.HaterDistanse {
			speed = common.HaterMaxSpeed
		} else if distance > common.HaterDistanse-20 && distance < common.HaterDistanse {
			if !enemy.BulletShotActive {
				h.HaterShot(i, enemy.HX, enemy.HY, playerX, playerY, world)
			}
			speed = 0
		} else if distance < common.HaterDistanse {
			if !enemy.BulletShotActive {
				h.HaterShot(i, enemy.HX, enemy.HY, playerX, playerY, world)
			}
			speed = float64(-common.HaterBaseSpeed)
		}

		if enemy.HatersCooldownActive {
			speed = speed * 5
		}

		if distance > 0.01 {
			enemy.HX += (dx / distance) * speed * dt
			enemy.HY += (dy / distance) * speed * dt
		}
		if common.Score%50 == 0 {
			common.HatersValwe++
			common.PatheticBaseSpeed += 50
			common.Score++
		}
		if CheckCollisionWithAttachment(enemy.HX, enemy.HY, common.HaterSize, attahmentX,
			attahmentY, float64(attahmentW), float64(attahmentH), attahmentAngle) {
			h.HatersMass = append(h.HatersMass[:i], h.HatersMass[i+1:]...)
			i--

			common.Score++

			common.PlayerHelth += 5

			player.ActivateBoost()

			continue
		}
		if CheckCollisionWithPlayer(enemy.HX, enemy.HY, common.HaterSize, world) {

			HatersCooldown(enemy)

			//ИГРОК ДОЛЖЕН ОТТАЛКИВАТЬ ОТ СЕБЯ ВРАГОВ

		}

	}
}

func (h *Haters) UpdateBulets(dt float64, world game.WorldView) {

	_, _, attahmentX, attahmentY, attahmentW, attahmentH, attahmentAngle := WhereThePlayer(world)

	for b := 0; b < len(h.Bullets); b++ {
		bulet := &h.Bullets[b]

		bulet.BX += bulet.VX * dt
		bulet.BY += bulet.VY * dt

		if bulet.BX < -50 || bulet.BX > common.ScreenWidth+50 ||
			bulet.BY < -50 || bulet.BY > common.ScreenWidth+50 {
			h.Bullets = append(h.Bullets[:b], h.Bullets[b+1:]...)
			b--
			continue
		}

		if CheckCollisionWithAttachment(bulet.BX, bulet.BY, common.HaitersBolletSize, attahmentX, attahmentY,
			attahmentW, attahmentH, attahmentAngle*math.Pi/180) {
			h.Bullets = append(h.Bullets[:b], h.Bullets[b+1:]...)
			b--
			continue
		}
		if CheckCollisionWithPlayer(bulet.BX, bulet.BY, common.HaitersBolletDamage, world) {
			common.PlayerHelth -= common.HaitersBolletDamage

			h.Bullets = append(h.Bullets[:b], h.Bullets[b+1:]...)
			b--
			continue
		}
	}
}

func (h *Haters) Update(worldView game.WorldView) bool {
	dt := 1.0 / 60.0

	if common.HatersValwe > common.MaxHatersValwe {
		common.HatersValwe = common.MaxHatersValwe
	}

	h.UpdateHaiters(dt, worldView)
	h.UpdateBulets(dt, worldView)
	return false
}

func (h *Haters) Draw(screen *ebiten.Image) {
}

func (h *Haters) Tag() string {
	return "hater"
}
