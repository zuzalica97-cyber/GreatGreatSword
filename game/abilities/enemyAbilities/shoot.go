package enemyabilities

import (
	"fmt"
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

var _ EnemyAbility = (*ShootAbility)(nil)
var _ hitboxes.HitBoxer = (*EnemyBullet)(nil)

// ============================================================
// ПУЛЯ
// ============================================================

type EnemyBullet struct {
	X, Y    float64
	VX, VY  float64
	Active  bool
	Damage  float64
	Speed   float64
	Size    float64
	Color   color.RGBA
	Index   int
	manager *hitboxes.CollisionManager
}

func (b *EnemyBullet) GetAABB() (posX, posY, halfW, halfH float64) {
	halfSize := b.Size / 2
	return b.X + halfSize, b.Y + halfSize, halfSize, halfSize
}

func (b *EnemyBullet) GetHitBoxID() string {
	return "enemyBullet_" + string(rune(b.Index))
}

func (b *EnemyBullet) IsActive() bool {
	return b.Active
}

func (b *EnemyBullet) IsStatic() bool {
	return false
}

func (b *EnemyBullet) ApplyPush(x, y float64) {}

func (b *EnemyBullet) OnCollision(other hitboxes.HitBoxer) {
	otherID := other.GetHitBoxID()
	switch otherID {
	case "blueSword", "playerLeg":
		b.Active = false
		fmt.Println(8)
		if otherID == "playerLeg" {
			common.PlayerHelth -= b.Damage
		}
	}
}

// ============================================================
// СПОСОБНОСТЬ ВЫСТРЕЛА
// ============================================================

type ShootAbility struct {
	// Параметры пули
	BulletSpeed  float64
	BulletDamage float64
	BulletSize   float64
	BulletColor  color.RGBA

	// Перезарядка
	CooldownMax float64
	Cooldown    float64

	// Состояние
	Active      bool
	ShotTimer   float64
	Bullets     []*EnemyBullet
	bulletIndex int
}

func NewShootAbility(bulletSpeed, bulletDamage, bulletSize, cooldown float64) *ShootAbility {
	return &ShootAbility{
		BulletSpeed:  bulletSpeed,
		BulletDamage: bulletDamage,
		BulletSize:   bulletSize,
		BulletColor:  color.RGBA{20, 20, 20, 255},
		CooldownMax:  cooldown,
		Bullets:      make([]*EnemyBullet, 0),
	}
}

func (s *ShootAbility) Name() string {
	return "Shoot"
}

// CreateBullet - создаёт пулю в сторону игрока
func (s *ShootAbility) CreateBullet(enemyX, enemyY, playerX, playerY float64, manager *hitboxes.CollisionManager) {
	dx := playerX - enemyX
	dy := playerY - enemyY
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0.01 {
		vx := (dx / distance) * s.BulletSpeed
		vy := (dy / distance) * s.BulletSpeed

		bullet := &EnemyBullet{
			X:      enemyX,
			Y:      enemyY,
			VX:     vx,
			VY:     vy,
			Active: true,
			Damage: s.BulletDamage,
			Speed:  s.BulletSpeed,
			Size:   s.BulletSize,
			Color:  s.BulletColor,
			Index:  s.bulletIndex,
		}
		s.bulletIndex++

		s.Bullets = append(s.Bullets, bullet)
		if manager != nil {
			manager.AddObject(bullet)
		}
	}
}

// Update - обновляет пули и перезарядку
func (s *ShootAbility) Update(enemy EnemyUser, dt float64, manager *hitboxes.CollisionManager) bool {
	// Обновляем кулдаун
	if s.Cooldown > 0 {
		s.Cooldown -= dt
	}

	// Обновляем пули
	for i := 0; i < len(s.Bullets); i++ {
		bullet := s.Bullets[i]
		if !bullet.Active {
			if manager != nil {
				manager.RemoveObject(bullet)
			}
			s.Bullets[i] = nil
			s.Bullets = append(s.Bullets[:i], s.Bullets[i+1:]...)
			i--
			continue
		}

		bullet.X += bullet.VX * dt
		bullet.Y += bullet.VY * dt

		// Проверка выхода за границы
		if bullet.X < -50 || bullet.X > common.RoomWidth+50 ||
			bullet.Y < -50 || bullet.Y > common.RoomHeight+50 {
			bullet.Active = false
		}
	}

	s.Active = false
	return false
}

// Activate - стреляет в сторону игрока
func (s *ShootAbility) Activate(enemy EnemyUser, world game.WorldView) bool {
	// Проверяем кулдаун
	if s.Cooldown > 0 {
		return false
	}

	// Получаем позицию игрока
	tx, ty := enemy.GetTarget()
	ex, ey := enemy.GetPosition()

	// Проверяем, что цель существует
	if tx == 0 && ty == 0 {
		return false
	}

	// Стреляем
	s.CreateBullet(ex, ey, tx, ty, nil)
	s.Cooldown = s.CooldownMax
	s.Active = true

	return true
}

func (s *ShootAbility) IsActive() bool {
	return s.Active
}

func (s *ShootAbility) GetCooldown() float64 {
	return s.Cooldown
}

func (s *ShootAbility) SetCooldown(value float64) {
	s.Cooldown = value
}

// Draw - рисует все пули
func (s *ShootAbility) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	for _, bullet := range s.Bullets {
		if !bullet.Active {
			continue
		}
		screenX := bullet.X - camera.X
		screenY := bullet.Y - camera.Y
		vector.FillRect(
			screen,
			float32(screenX),
			float32(screenY),
			float32(bullet.Size),
			float32(bullet.Size),
			bullet.Color,
			true,
		)
	}
}
