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
	"github.com/hajimehoshi/ebiten/v2/vector" // ← импорт библиотеки
	"github.com/setanarut/kamera/v2"
)

var _ game.Entity = (*Pathetic)(nil)
var _ hitboxes.HitBoxer = (*OnePath)(nil)

type Pathetic struct {
	Paths []*OnePath
}

type OnePath struct {
	Size                   int
	PX, PY                 float64
	Helth                  int
	Damage                 int
	Speed                  float64
	MaxSpeed               float64
	Distance               float64
	PathericCooldownActive bool
	PathericCooldownTimer  float64
	Active                 bool
	Texture                *ebiten.Image
	Color                  color.RGBA

	// Поля для рывка
	DashActive      bool
	DashDistance    float64
	DashTimer       float64
	DashTimerMax    float64
	DashTargetX     float64
	DashTargetY     float64
	DashDirX        float64
	DashDirY        float64
	DashSpeed       float64
	DashSpeedMax    float64
	DashCooldown    float64
	DashCooldownMax float64
}

func NewPathetic() *Pathetic {
	p := &Pathetic{}
	return p
}

func (p *Pathetic) SpawnPathetic(x, y float64, manager *hitboxes.CollisionManager) {
	enemy := &OnePath{
		Size:     65,
		PX:       x,
		PY:       y,
		Active:   true,
		Helth:    50,
		Damage:   5,
		Speed:    150,
		MaxSpeed: 200,
		Distance: 400,
		Color:    color.RGBA{150, 150, 150, 255},

		DashActive:      false,
		DashDistance:    80,
		DashTimer:       0.9,
		DashTimerMax:    0.9,
		DashSpeed:       500,
		DashSpeedMax:    500,
		DashCooldown:    0,
		DashCooldownMax: 2.0,
	}

	p.Paths = append(p.Paths, enemy)
	manager.AddObject(enemy)
}

// PerformDash - выполняет рывок с ФИКСИРОВАННЫМ направлением
func (enemy *OnePath) PerformDash(targetX, targetY float64, dashSpeed, dashDuration float64) {
	dx := targetX - enemy.PX
	dy := targetY - enemy.PY
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0.01 {
		enemy.DashDirX = dx / distance
		enemy.DashDirY = dy / distance
		enemy.DashSpeed = dashSpeed
		enemy.DashActive = true
		enemy.DashTimer = dashDuration
	}
}

// UpdateDash - обновляет состояние рывка
func (enemy *OnePath) UpdateDash(dt float64) bool {
	if !enemy.DashActive {
		return false
	}

	enemy.DashTimer -= dt

	enemy.PX += enemy.DashDirX * enemy.DashSpeed * dt
	enemy.PY += enemy.DashDirY * enemy.DashSpeed * dt

	if enemy.PX < 0 || enemy.PX > common.RoomWidth ||
		enemy.PY < 0 || enemy.PY > common.RoomHeight {
		enemy.DashActive = false
		enemy.DashCooldown = enemy.DashCooldownMax
		return true
	}

	if enemy.DashTimer <= 0 {
		enemy.DashActive = false
		enemy.DashSpeed = 0
		enemy.DashCooldown = enemy.DashCooldownMax
		return true
	}

	return false
}

// ===== РЕАЛИЗАЦИЯ INTERFACE HitBoxer =====

// GetAABB - возвращает AABB для проверки коллизий
func (e *OnePath) GetAABB() (posX, posY, halfW, halfH float64) {
	halfSize := float64(e.Size) / 2
	return e.PX + halfSize, e.PY + halfSize, halfSize, halfSize
}

// GetHitBoxID - возвращает уникальный ID
func (e *OnePath) GetHitBoxID() string {
	return fmt.Sprintf("enemy_%p", e) // используем адрес как ID
}

// IsActive - проверяет, активен ли враг
func (e *OnePath) IsActive() bool {
	return e.Active
}

// IsStatic - враг не статичен (двигается)
func (e *OnePath) IsStatic() bool {
	return false
}

// ApplyPush - применяет силу отталкивания
func (e *OnePath) ApplyPush(x, y float64) {
	switch e.DashActive {
	case false:
		e.PX += x
		e.PY += y
	}
}

// OnCollision - вызывается при столкновении с другим объектом
func (e *OnePath) OnCollision(other hitboxes.HitBoxer) {
}

func (p *Pathetic) UpatePathetics(dt float64, worldView game.WorldView) {

	playerX, playerY, attahmentX, attahmentY, attahmentW, attahmentH, attahmentAngle := WhereThePlayer(worldView)

	for i := 0; i < len(p.Paths); i++ {
		enemy := p.Paths[i]
		if !enemy.Active {
			continue
		}

		if enemy.PathericCooldownActive {
			enemy.PathericCooldownTimer -= dt
			if enemy.PathericCooldownTimer <= 0 {
				enemy.PathericCooldownActive = false
			}
		}

		// ===== КУЛДАУН РЫВКА (ВСЕГДА) =====
		if enemy.DashCooldown > 0 {
			enemy.DashCooldown -= dt
			if enemy.DashCooldown <= 0 {
				enemy.DashSpeed = enemy.DashSpeedMax
				enemy.DashTimer = enemy.DashTimerMax
			}
		}

		// ===== ОБНОВЛЕНИЕ РЫВКА =====
		if enemy.DashActive {
			enemy.UpdateDash(dt)
			if !enemy.DashActive {
				continue
			}
		}

		if !enemy.DashActive {
			dx := playerX - enemy.PX
			dy := playerY - enemy.PY
			distance := math.Sqrt(dx*dx + dy*dy)

			var speed float64

			if distance > enemy.Distance {
				speed = enemy.Speed
			} else {
				t := 1.0 - distance/enemy.Distance
				speed = enemy.Speed + t*float64(enemy.MaxSpeed-enemy.Speed) // решить пробемму с залипанием и рванным движением. после добавить коллизи к остальным обьектам а также добавить механику веса и отпрежинивание при столканвении.
			}

			// ===== АКТИВАЦИЯ РЫВКА =====
			if distance < enemy.Distance && !enemy.DashActive && enemy.DashCooldown <= 0 {
				enemy.PerformDash(playerX, playerY, enemy.DashSpeed, enemy.DashTimer)
				continue
			}

			if enemy.PathericCooldownActive {
				speed = -100
			}

			if distance > 0.01 {
				enemy.PX += (dx / distance) * speed * dt
				enemy.PY += (dy / distance) * speed * dt
			}
		}

		// ===== КОЛЛИЗИИ С МЕЧОМ =====
		if CheckCollisionWithAttachment(enemy.PX, enemy.PY, float64(enemy.Size), attahmentX,
			attahmentY, float64(attahmentW), float64(attahmentH), attahmentAngle) {

			if !enemy.PathericCooldownActive || enemy.PathericCooldownTimer <= 0.5 {
				enemy.Helth -= common.PlayerDamage
				PatheticCooldown(enemy)
			}

			if enemy.Helth <= 0 {
				p.Paths = append(p.Paths[:i], p.Paths[i+1:]...)
				i--
				common.Score++
				player.ActivateBoost()
			}
			continue
		}

		// ===== КОЛЛИЗИИ С ИГРОКОМ =====
		if CheckCollisionWithPlayer(enemy.PX, enemy.PY, float64(enemy.Size), worldView) {

			if !enemy.PathericCooldownActive {
				for _, pleg := range worldView.SearchEntities("playerLeg") {
					p := pleg.(*player.PlayerLeg)
					p.DamagePlayer(enemy.Damage)
				}
			}

			PatheticCooldown(enemy)

			if !enemy.DashActive {
				dx := playerX - enemy.PX
				dy := playerY - enemy.PY
				distance := math.Sqrt(dx*dx + dy*dy)
				if distance > 0.01 {
					enemy.PX += (dx / distance) * 50 * dt
					enemy.PY += (dy / distance) * 50 * dt
				}
			}
		}
	}
}

func (p *Pathetic) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool { //ДЗ нужно сделать проверку коллизей общеё для вех врагов также нужен интерфейс типа враг с функциями чтобы с ними взаимодействовать. после нужено добавить врагов с оружиями и добавить эффекты для игрока и варогов.
	// И нужно подпаравит деш чтобы у нас был не 1 а 3 и визальные эффекты к способностям.
	dt := 1.0 / 60.0

	if len(p.Paths) < 8 {
		x, y := RangomSpawnInWall(50)
		p.SpawnPathetic(x, y, manager)
	}

	if common.Valwe > common.MaxValwe {
		common.Valwe = common.MaxValwe
	}

	p.UpatePathetics(dt, worldView)

	return false
}

func (p *Pathetic) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	for _, enemy := range p.Paths {
		screenX := enemy.PX - camera.X
		screenY := enemy.PY - camera.Y

		Color := enemy.Color
		if enemy.PathericCooldownActive {
			Color = color.RGBA{120, 120, 120, 255}
		}

		vector.FillRect(
			screen,
			float32(screenX),
			float32(screenY),
			float32(enemy.Size),
			float32(enemy.Size),
			Color,
			true,
		)
	}
}

func (p *Pathetic) Tag() string {
	return "pathetic"
}
