package enemy

import (
	"great-sword/game"
	enemyabilities "great-sword/game/abilities/enemyAbilities"
	"great-sword/game/hitboxes"
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
			65,   // size
			40,   // helth
			5,    // damage
			100,  // speed
			300,  // maxSpeed
			0.7,  // чем ближе к 1, тем дольше скользит
			70.0, // как быстро разгоняется
			color.RGBA{150, 150, 150, 255},
			2,
			0.6,
			"pathetic",
			180.0, // 360 градусов в секунду
			720.0, // ускорение
			540.0, // замедление
			0.15,  // плавность
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

	if len(p.Paths) < 5 {
		x, y := RangomSpawnInWall(50)
		p.SpawnPathetic(x, y, manager)
	}

	for i := 0; i < len(p.Paths); i++ {
		enemy := p.Paths[i]

		// Проверка смерти
		var died bool
		p.Paths, died = DeathScan(manager, enemy.BaseEnemy, p.Paths, i)
		if died {
			i-- // корректируем индекс после удаления
			continue
		}

		enemy.StUpdateCoolDown(dt, worldView, manager)

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

		if enemy.CooldownActive {
			enemy.CurrentSpeed = -enemy.CurrentSpeed / 3
		}

		enemy.StMovment(dt)
	}

	return false
}

func (b *Pathetic) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	for _, b := range b.Paths {
		b.StDraw(screen, camera)
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
