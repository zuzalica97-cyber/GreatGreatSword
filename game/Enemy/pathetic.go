package enemy

import (
	"fmt"
	"great-sword/game"
	enemyabilities "great-sword/game/abilities/enemyAbilities"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
	"great-sword/game/player"
	"great-sword/game/weapons"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

var _ game.Entity = (*Pathetic)(nil)
var _ hitboxes.HitBoxer = (*OnePath)(nil)
var _ enemyabilities.EnemyUser = (*OnePath)(nil)
var _ hitboxes.LetterReceiver = (*OnePath)(nil)
var _ game.Enemy = (*OnePath)(nil)

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
			500, // helth
			5,   // damage
			100, // speed
			300, // maxSpeed
			color.RGBA{150, 150, 150, 255},
			3,
			0.6,
			"pathetic",
		),
	}

	// Добавляем способности
	enemy.Abilities = []enemyabilities.EnemyAbility{
		enemyabilities.NewChaseAbility(enemy.Speed, enemy.MaxSpeed, 600, 0.01),
		enemyabilities.NewDashAbility(250, 450, 2, 0.5),
	}

	enemy.SetWeapon(weapons.NewlitleKnife(manager, enemy))

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

// ============================================================
// ОБНОВЛЕНИЕ
// ============================================================

func (p *Pathetic) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {
	dt := 1.0 / 60.0

	playerX, playerY := getPlayerPosition(worldView) //ДЗ делать мечи у партивников. и делать их крутищихся

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
				enemy.RemoveWeaponFromCollision(manager)
			}
			p.Paths[i] = nil
			p.Paths = append(p.Paths[:i], p.Paths[i+1:]...)
			i--
			common.Score++
			player.ActivateBoost()
			continue
		}

		// ===== ВРАЩЕНИЕ ВПРАВО =====
		// Вращаемся вправо (1 = вправо, -1 = влево)
		enemy.Rotation.UpdateRotation(1.0, dt)

		// Обновляем оружие (если есть)
		if enemy.weapon != nil {
			enemy.weapon.Update(worldView, manager)
		}

		// === ОБНОВЛЕНИЕ КУЛДАУНА ===
		enemy.UpdateCooldown(dt)

		for _, letter := range enemy.Letters {
			letter.UpdateCoolDown(dt)
		}

		// === УСТАНОВКА ЦЕЛИ ===
		enemy.SetTarget(playerX, playerY)

		enemy.CurrentSpeed = enemy.Speed

		// === ОБНОВЛЕНИЕ СПОСОБНОСТЕЙ ===
		for _, ability := range enemy.Abilities {
			ability.Activate(enemy, worldView)
			ability.Update(enemy, dt, manager)
		}

		// === ОБНОВЛЕНИЕ ЭФФЕКТОВ ===
		enemy.UpdateEffects(dt)

		if len(enemy.Effects) > 0 {
			fmt.Println(len(enemy.Effects))
		}

		if enemy.CooldownActive {
			enemy.CurrentSpeed = -enemy.CurrentSpeed / 3
		}

		friction := 0.7      // чем ближе к 1, тем дольше скользит
		acceleration := 70.0 // как быстро разгоняется

		newSpeed, newDirX, newDirY := enemy.EnemySlideMovmentFunc(friction, acceleration, dt)

		enemy.CurrentSpeed = newSpeed
		enemy.SetDirection(newDirX, newDirY)

		// === ДВИЖЕНИЕ ===
		newX, newY := MoveEnemyToTareget(enemy.BaseEnemy, dt)

		enemy.SetPosition(newX, newY)

	}

	return false
}

func (b *Pathetic) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	for _, b := range b.Paths {
		// Центр игрока в мировых координатах
		centerX := b.X + float64(b.Size)/2
		centerY := b.Y + float64(b.Size)/2

		// Экранные координаты с учётом камеры
		screenX := centerX - camera.X
		screenY := centerY - camera.Y

		// 1. Рисуем повёрнутый квадрат
		DrawRotatedRect(
			screen,
			screenX,
			screenY,
			float64(b.Size),
			float64(b.Size),
			b.Rotation.Angle, // ← угол поворота
			b.Color,
		)

		for _, effect := range b.Effects {
			effect.Draw(screen, camera, b)
		}

		if b.weapon != nil {
			b.weapon.Draw(screen, camera)
		}
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
