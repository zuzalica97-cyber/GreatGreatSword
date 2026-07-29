package enemy

import (
	"great-sword/game"
	enemyabilities "great-sword/game/abilities/enemyAbilities"
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

type Haters struct {
	HatersMass   []*Haits
	frameCounter float64
}

// ============================================================
// Haits - ВРАГ (использует BaseEnemy)
// ============================================================

type Haits struct {
	*BaseEnemy // ← ВСТРАИВАНИЕ!

	// Специфичные поля
	ShotActive bool
	ShotTimer  float64
	Cooldown   float64
	Distanse   float64
	Abilities  []enemyabilities.EnemyAbility
}

func NawHaters() *Haters {
	h := &Haters{
		frameCounter: 1.0,
	}
	return h
}

// ============================================================
// СОЗДАНИЕ ВРАГА
// ============================================================

func (h *Haters) SpawnHiters(x, y float64, manager *hitboxes.CollisionManager) {
	enemy := &Haits{
		BaseEnemy: NewBaseEnemy(
			x, y,
			50,  // size
			35,  // health
			5,   // damage
			150, // baseSpeed
			300, // maxSpeed
			color.RGBA{180, 150, 150, 255},
			"hater",
		),
		Distanse: 500,
		Cooldown: 1.0,
	}

	// Добавляем способности
	enemy.Abilities = []enemyabilities.EnemyAbility{

		enemyabilities.NewShootAbility(200, 10, 15, 1.0), // ← способность выстрела
		enemyabilities.NewFleeAbility(300, 200, 500, 50),
	}

	h.HatersMass = append(h.HatersMass, enemy)
	if manager != nil {
		manager.AddObject(enemy)
	}
}

// ============================================================
// OnCollision - обработка столкновений
// ============================================================

func (e *Haits) OnCollision(other hitboxes.HitBoxer) {
	otherID := other.GetHitBoxID()
	switch otherID {
	case "blueSword":
		e.TakeDamage(common.PlayerDamage)
	case "playerLeg":
		HatersCooldown(e)
	}
}

// ============================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ============================================================

func HatersCooldown(enemy *Haits) {
	enemy.CooldownActive = true
	enemy.CooldownTimer = 2.0
}

// ============================================================
// ЛОГИКА ПРОТИВНИКА
// ============================================================

func (h *Haters) UpdateHaiters(dt float64, world game.WorldView, manager *hitboxes.CollisionManager) {
	var playerX, playerY float64

	if len(h.HatersMass) == 0 {
		return
	}

	// Получаем позицию игрока
	for _, pleg := range world.SearchEntities("playerLeg") {
		p, ok := pleg.(game.PlayerLegInter)
		if !ok {
			continue
		}
		playerX, playerY = p.GetPosition()
	}

	for i := 0; i < len(h.HatersMass); i++ {
		enemy := h.HatersMass[i]

		// === ПРОВЕРКА СМЕРТИ ===
		if !enemy.IsActive() || enemy.GetHealth() <= 0 {
			if manager != nil {
				manager.RemoveObject(enemy)
			}
			h.HatersMass[i] = nil
			h.HatersMass = append(h.HatersMass[:i], h.HatersMass[i+1:]...)
			i--
			common.Score++
			player.ActivateBoost()
			continue
		}

		// === ОБНОВЛЕНИЕ КУЛДАУНА ===
		enemy.UpdateCooldown(dt)

		// === УСТАНОВКА ЦЕЛИ ===
		enemy.SetTarget(playerX, playerY)

		// === АКТИВАЦИЯ И ОБНОВЛЕНИЕ СПОСОБНОСТЕЙ ===
		for _, ability := range enemy.Abilities {
			ability.Update(enemy, dt, manager)
		}

		// === РАСЧЁТ СКОРОСТИ ===
		XPos, YPos := enemy.GetPosition()
		dx := playerX - XPos
		dy := playerY - YPos
		distance := math.Sqrt(dx*dx + dy*dy)

		enemy.CurrentSpeed = enemy.Speed

		if distance < enemy.Distanse {
			for _, ability := range enemy.Abilities {
				switch ability.Name() {
				case "Shoot":
					ability.Activate(enemy, world)
					continue
				}
			}
		}

		if enemy.CooldownActive {
			enemy.CurrentSpeed = enemy.CurrentSpeed * 5
		}

		// === ДВИЖЕНИЕ ===
		if enemy.GetTargetDistance() > 0.01 && enemy.GetSpeed() != 0 {
			newX, newY := MoveEnemyTowardsPlayer(
				enemy.X, enemy.Y,
				playerX, playerY,
				enemy.GetSpeed(),
				dt,
			)
			enemy.SetPosition(newX, newY)
		}
	}
}

// ============================================================
// ОБНОВЛЕНИЕ
// ============================================================

func (h *Haters) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {
	dt := 1.0 / 60.0

	// Очистка слайсов
	if h.frameCounter > 0 {
		h.frameCounter -= dt
		if h.frameCounter <= 0 {
			h.frameCounter = 1.0
			h.HatersMass = hitboxes.CompactSlice(h.HatersMass)
		}
	}

	// Спавн врагов
	if len(h.HatersMass) < 1 {
		x, y := RangomSpawnInWall(50)
		h.SpawnHiters(x, y, manager)
	}

	h.UpdateHaiters(dt, worldView, manager)

	// Обновляем пули из способностей
	for _, enemy := range h.HatersMass {
		for _, ability := range enemy.Abilities {
			if ability.Name() == "Shoot" {
				ability.Update(enemy, dt, manager)
			}
		}
	}

	return false
}

// ============================================================
// ОТРИСОВКА
// ============================================================

func (h *Haters) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	// Рисуем врагов
	for _, enemy := range h.HatersMass {
		screenX := enemy.X - camera.X
		screenY := enemy.Y - camera.Y

		Color := enemy.Color
		if enemy.CooldownActive {
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

	// Рисуем пули из способностей
	for _, enemy := range h.HatersMass {
		for _, ability := range enemy.Abilities {
			if ability.Name() == "Shoot" {
				if shoot, ok := ability.(*enemyabilities.ShootAbility); ok {
					shoot.Draw(screen, camera)
				}
			}
		}
	}
}

// ============================================================
// ИНТЕРФЕЙС game.Entity
// ============================================================

func (h *Haters) Tag() string {
	return "hater"
}

func (h *Haters) IsActive() bool {
	return true
}
