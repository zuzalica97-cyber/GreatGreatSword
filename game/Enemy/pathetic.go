package enemy

import (
	"great-sword/game"
	enemyabilities "great-sword/game/abilities/enemyAbilities"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
	"great-sword/game/weapons"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

var _ game.Enemy = (*OnePath)(nil)
var _ game.Entity = (*OnePath)(nil)
var _ hitboxes.HitBoxer = (*OnePath)(nil)
var _ enemyabilities.EnemyUser = (*OnePath)(nil)
var _ hitboxes.LetterReceiver = (*OnePath)(nil)

// ============================================================
// ОСНОВНАЯ СТРУКТУРА ВРАГА
// ============================================================

type OnePath struct {
	*BaseEnemy

	PathericCooldownActive bool
	PathericCooldownTimer  float64

	Abilities []enemyabilities.EnemyAbility
}

// ============================================================
// СОЗДАНИЕ ВРАГА
// ============================================================

func SpawnPathetic(x, y float64, manager *hitboxes.CollisionManager) *OnePath {
	enemy := &OnePath{
		BaseEnemy: NewBaseEnemy(
			x, y,
			65,   // size
			4000, // health
			5,    // damage
			100,  // speed
			300,  // maxSpeed
			0.7,  // трение
			70.0, // ускорение
			color.RGBA{150, 150, 150, 255},
			2,
			0.6,
			"pathetic",
			180.0,
			720.0,
			540.0,
			0.15,
		),
	}

	enemy.Abilities = []enemyabilities.EnemyAbility{
		enemyabilities.NewChaseAbility(enemy.Speed, enemy.MaxSpeed, 600, 0.01),
		enemyabilities.NewDashAbility(250, 450, 2, 0.5),
	}

	enemy.AddWeapon(weapons.NewHammer(enemy, "left", manager))

	if manager != nil {
		manager.AddObject(enemy)
	}
	return enemy
}

// ============================================================
// ОБНОВЛЕНИЕ
// ============================================================

func (enemy *OnePath) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {
	dt := 1.0 / 60.0

	playerX, playerY := getPlayerPosition(worldView)

	// === УСТАНОВКА ЦЕЛИ ===
	enemy.SetTarget(playerX, playerY)

	// === ОБНОВЛЕНИЕ СПОСОБНОСТЕЙ ===
	enemy.CurrentSpeed = enemy.Speed

	for _, ability := range enemy.Abilities {
		ability.Activate(enemy, worldView)
		ability.Update(enemy, dt, manager)
	}

	// === ОБНОВЛЕНИЕ ЭФФЕКТОВ ===
	enemy.UpdateEffects(dt)

	enemy.StMovment(dt)

	// === КУЛДАУН (ОТТАЛКИВАНИЕ) ===
	if enemy.CooldownActive {
		enemy.CurrentSpeed = -enemy.CurrentSpeed / 3
	}

	// === ГРАНИЦЫ КОМНАТЫ ===
	if enemy.X < 0 {
		enemy.X = 0
		enemy.CurrentSpeed = -enemy.CurrentSpeed * 0.5
	}
	if enemy.X > common.RoomWidth-float64(enemy.Size) {
		enemy.X = common.RoomWidth - float64(enemy.Size)
		enemy.CurrentSpeed = -enemy.CurrentSpeed * 0.5
	}
	if enemy.Y < 0 {
		enemy.Y = 0
		enemy.CurrentSpeed = -enemy.CurrentSpeed * 0.5
	}
	if enemy.Y > common.RoomHeight-float64(enemy.Size) {
		enemy.Y = common.RoomHeight - float64(enemy.Size)
		enemy.CurrentSpeed = -enemy.CurrentSpeed * 0.5
	}

	return false
}

// ============================================================
// ОТРИСОВКА
// ============================================================

func (b *OnePath) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	b.StDraw(screen, camera)
}

// ============================================================
// ИНТЕРФЕЙС game.Entity
// ============================================================

func (b *OnePath) Tag() string {
	return b.TagName
}

func (b *OnePath) IsActive() bool {
	return b.Active
}
