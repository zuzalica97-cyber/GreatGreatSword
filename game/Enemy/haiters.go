package enemy

import (
	"fmt"
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
	"great-sword/game/player"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

var _ game.Entity = (*Haters)(nil)
var _ hitboxes.HitBoxer = (*Haits)(nil)
var _ hitboxes.HitBoxer = (*HaitersBullet)(nil)

type Haters struct {
	HatersMass []*Haits
	Bullets    []*HaitersBullet
}

type Haits struct {
	HatersCooldownActive bool
	HatersCooldownTimer  float64
	BulletShotActive     bool
	BulletShotTimer      float64
	HX, HY               float64
	Active               bool
	Texture              *ebiten.Image
	Color                color.RGBA
}
type HaitersBullet struct {
	BX, BY    float64
	VX, VY    float64
	BulActive bool
	Texture   *ebiten.Image
	Color     color.RGBA
}

func NawHaters() *Haters {
	h := &Haters{}
	return h
}

func (h *Haters) SpawnHiters(x, y float64, manager *hitboxes.CollisionManager) {

	enemy := &Haits{
		HX:     x,
		HY:     y,
		Active: true,
		Color:  color.RGBA{180, 150, 150, 255},
	}

	h.HatersMass = append(h.HatersMass, enemy)
	manager.AddObject(enemy)
}

func (h *Haters) CreateBullet(enemyX, enemyY, playerX, playerY float64, world game.WorldView, manager *hitboxes.CollisionManager) {

	dx := playerX - enemyX
	dy := playerY - enemyY
	diststance := math.Sqrt(dx*dx + dy*dy)

	if diststance > 0.01 {
		vx := (dx / diststance) * common.HaitersBolletSpeed
		vy := (dy / diststance) * common.HaitersBolletSpeed

		bullet := &HaitersBullet{
			BX:        enemyX,
			BY:        enemyY,
			VX:        vx,
			VY:        vy,
			BulActive: true,
		}

		h.Bullets = append(h.Bullets, bullet)
		manager.AddObject(bullet)
	}
}

// ============================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА hitboxes.HitBoxer ДЛЯ Haits
// ============================================================

// GetAABB - возвращает AABB для проверки коллизий
func (e *Haits) GetAABB() (posX, posY, halfW, halfH float64) {
	halfSize := common.HaterSize / 2
	return e.HX + halfSize, e.HY + halfSize, halfSize, halfSize
}

// GetHitBoxID - возвращает уникальный ID
func (e *Haits) GetHitBoxID() string {
	return fmt.Sprintf("hater_%p", e)
}

// IsActive - проверяет, активен ли враг
func (e *Haits) IsActive() bool {
	return e.Active
}

// IsStatic - враг не статичен
func (e *Haits) IsStatic() bool {
	return false
}

// ApplyPush - применяет силу отталкивания
func (e *Haits) ApplyPush(x, y float64) {
	e.HX += x
	e.HY += y
}

// OnCollision - вызывается при столкновении с другим объектом
func (e *Haits) OnCollision(other hitboxes.HitBoxer) {
	// Логика при столкновении (можно расширить)
}

// ============================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА hitboxes.HitBoxer ДЛЯ HaitersBullet
// ============================================================

// GetAABB - возвращает AABB для проверки коллизий
func (b *HaitersBullet) GetAABB() (posX, posY, halfW, halfH float64) {
	halfSize := common.HaitersBolletSize / 2
	return b.BX + halfSize, b.BY + halfSize, halfSize, halfSize
}

// GetHitBoxID - возвращает уникальный ID
func (b *HaitersBullet) GetHitBoxID() string {
	return fmt.Sprintf("bullet_%p", b)
}

// IsActive - проверяет, активна ли пуля
func (b *HaitersBullet) IsActive() bool {
	return b.BulActive
}

// IsStatic - пуля не статична
func (b *HaitersBullet) IsStatic() bool {
	return false
}

// ApplyPush - применяет силу отталкивания (пуля не отталкивается)
func (b *HaitersBullet) ApplyPush(x, y float64) {
	// Пуля не отталкивается, она либо летит, либо удаляется
	// Можно оставить пустым
}

// OnCollision - вызывается при столкновении с другим объектом
func (b *HaitersBullet) OnCollision(other hitboxes.HitBoxer) {
	// При столкновении пуля удаляется (логика в UpdateBulets)
}

//ЛОГИКА ПРАТИВНИКА

func (h *Haters) HaterShot(number int, enemyX, enemyY, playerX, playerY float64, world game.WorldView, manager *hitboxes.CollisionManager) {
	enemy := h.HatersMass[number]

	enemy.BulletShotActive = true
	enemy.BulletShotTimer = common.HaitersBolletCooldown
	h.CreateBullet(enemyX, enemyY, playerX, playerY, world, manager)
}

func (h *Haters) UpdateHaiters(dt float64, world game.WorldView, manager *hitboxes.CollisionManager) {

	playerX, playerY, attahmentX, attahmentY, attahmentW, attahmentH, attahmentAngle := WhereThePlayer(world)

	for i := 0; i < len(h.HatersMass); i++ {
		enemy := h.HatersMass[i]
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
				h.HaterShot(i, enemy.HX, enemy.HY, playerX, playerY, world, manager)
			}
			speed = 0
		} else if distance < common.HaterDistanse {
			if !enemy.BulletShotActive {
				h.HaterShot(i, enemy.HX, enemy.HY, playerX, playerY, world, manager)
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
		bulet := h.Bullets[b]

		bulet.BX += bulet.VX * dt
		bulet.BY += bulet.VY * dt

		if bulet.BX < -50 || bulet.BX > common.RoomWidth+50 ||
			bulet.BY < -50 || bulet.BY > common.RoomHeight+50 {
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

func (h *Haters) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {
	dt := 1.0 / 60.0

	if len(h.HatersMass) < 5 {
		x, y := RangomSpawnInWall(common.HaterSize)
		h.SpawnHiters(x, y, manager)
	}

	h.UpdateHaiters(dt, worldView, manager)
	h.UpdateBulets(dt, worldView)
	return false
}

func (h *Haters) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	for _, enemy := range h.HatersMass {

		screenX := enemy.HX - camera.X
		screenY := enemy.HY - camera.Y

		vector.FillRect(
			screen,
			float32(screenX),
			float32(screenY),
			common.HaterSize,
			common.HaterSize,
			enemy.Color,
			true,
		)
	}
	for _, bullet := range h.Bullets {

		screenX := bullet.BX - camera.X
		screenY := bullet.BY - camera.Y
		vector.FillRect(
			screen,
			float32(screenX),
			float32(screenY),
			common.HaitersBolletSize,
			common.HaitersBolletSize,
			color.RGBA{20, 20, 20, 255},
			true,
		)
	}
}

func (h *Haters) Tag() string {
	return "hater"
}
