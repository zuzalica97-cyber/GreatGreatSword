package enemy

import (
	"great-sword/game"
	enemyabilities "great-sword/game/abilities/enemyAbilities"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
	"great-sword/game/player"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

var _ game.Entity = (*Pathetic)(nil)
var _ hitboxes.HitBoxer = (*OnePath)(nil)
var _ enemyabilities.EnemyUser = (*OnePath)(nil)

// ============================================================
// ОСНОВНАЯ СТРУКТУРА ВРАГА
// ============================================================

type Pathetic struct {
	Paths []*OnePath
}

type OnePath struct {
	*BaseEnemy // ← ВСТРАИВАНИЕ! Все методы BaseEnemy доступны

	// Специфичные поля (только то, чего нет в BaseEnemy)
	PathericCooldownActive bool
	PathericCooldownTimer  float64

	// Способности
	Abilities []enemyabilities.EnemyAbility

	// Кэш для Target (позиция игрока) - уже есть в BaseEnemy
}

// ============================================================
// КОНСТРУКТОР
// ============================================================

func NewPathetic() *Pathetic {
	return &Pathetic{
		Paths: make([]*OnePath, 0),
	}
}

// ============================================================
// СОЗДАНИЕ ВРАГА
// ============================================================

func (p *Pathetic) SpawnPathetic(x, y float64, manager *hitboxes.CollisionManager) {
	enemy := &OnePath{
		BaseEnemy: NewBaseEnemy(
			x, y,
			65,  // size
			50,  // health
			5,   // damage
			100, // baseSpeed
			300, // maxSpeed
			color.RGBA{150, 150, 150, 255},
			"pathetic",
		),
	}

	// Добавляем способности
	enemy.Abilities = []enemyabilities.EnemyAbility{
		enemyabilities.NewChaseAbility(enemy.Speed, enemy.MaxSpeed, 600, 250),
		enemyabilities.NewDashAbility(250, 500, 2, 0.3),
	}

	p.Paths = append(p.Paths, enemy)
	if manager != nil {
		manager.AddObject(enemy)
	}
}

// PatheticCooldown - специфичная для этого врага
func PatheticCooldown(enemy *OnePath) {
	enemy.CooldownActive = true
	enemy.CooldownTimer = 2.0
}

func (e *OnePath) OnCollision(other hitboxes.HitBoxer) {
	otherID := other.GetHitBoxID()
	switch otherID {
	case "blueSword":
		e.TakeDamage(common.PlayerDamage)
	case "playerLeg":
		if !e.CooldownActive {
			common.PlayerHelth -= float64(e.Damage)
			PatheticCooldown(e)
		} //ДЗ нужно сделать нармальное нанесение урона. сделать нормальный фиал Haters и для него способности быть на дистации и гинерировать пули
	}
}

// ============================================================
// ОБНОВЛЕНИЕ
// ============================================================

func (p *Pathetic) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {
	dt := 1.0 / 60.0

	playerX, playerY := getPlayerPosition(worldView)

	if len(p.Paths) < 1 {
		x, y := RangomSpawnInWall(50)
		p.SpawnPathetic(x, y, manager)
	}

	for i := 0; i < len(p.Paths); i++ {
		enemy := p.Paths[i]

		// === ПРОВЕРКА СМЕРТИ ===
		if !enemy.IsActive() || enemy.GetHealth() <= 0 {
			if manager != nil {
				manager.RemoveObject(enemy)
			}
			p.Paths[i] = nil
			p.Paths = append(p.Paths[:i], p.Paths[i+1:]...)
			i--
			common.Score++
			player.ActivateBoost()
			continue
		}

		// === ОБНОВЛЕНИЕ КУЛДАУНА ===
		enemy.UpdateCooldown(dt)

		// === УСТАНОВКА ЦЕЛИ ===
		enemy.SetTarget(playerX, playerY)

		enemy.CurrentSpeed = enemy.Speed

		// === ОБНОВЛЕНИЕ СПОСОБНОСТЕЙ ===
		for _, ability := range enemy.Abilities {
			ability.Activate(enemy, worldView)
			ability.Update(enemy, dt, manager)
		}

		if enemy.CooldownActive {
			enemy.CurrentSpeed = -enemy.CurrentSpeed / 3
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

	return false
}

// ============================================================
// ОТРИСОВКА
// ============================================================

func (p *Pathetic) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	for _, enemy := range p.Paths {
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
}

// ============================================================
// ИНТЕРФЕЙС game.Entity
// ============================================================

func (p *Pathetic) Tag() string {
	return "pathetic"
}

func (p *Pathetic) IsActive() bool {
	return true
}
