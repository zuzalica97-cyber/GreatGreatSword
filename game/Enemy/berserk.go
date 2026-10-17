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

var _ game.Enemy = (*OneBerserk)(nil)
var _ game.Entity = (*OneBerserk)(nil)
var _ hitboxes.HitBoxer = (*OneBerserk)(nil)
var _ enemyabilities.EnemyUser = (*OneBerserk)(nil)

// ============================================================
// ОСНОВНАЯ СТРУКТУРА ВРАГА
// ============================================================

type OneBerserk struct {
	*BaseEnemy
	Abilities []enemyabilities.EnemyAbility
	oldSpeed  float64
}

// ============================================================
// КОНСТРУКТОР
// ============================================================

// ============================================================
// СОЗДАНИЕ ВРАГА
// ============================================================

func SpawnBreserk(x, y float64, manager *hitboxes.CollisionManager) *OneBerserk {
	enemy := &OneBerserk{
		BaseEnemy: NewBaseEnemy(
			x, y,
			150,                         // size
			450,                         // health
			10,                          // damage
			200,                         // baseSpeed
			400,                         // maxSpeed
			0.9999,                      // чем ближе к 1, тем дольше скользит
			50.0,                        // как быстро разгоняется
			color.RGBA{90, 90, 90, 255}, //ДЗ нужно поравить силу ооталкивания чтобы тяжёлые обьекты легко отталкивали лёгкие а то аура мешает
			50,
			1,
			"berserk",
			120.0, // 360 градусов в секунду
			500.0, // ускорение
			700.0, // замедление
			0.3,   // плавность
		),
		oldSpeed: 0,
	}

	enemy.Abilities = append(enemy.Abilities,
		enemyabilities.NewDashAbility(800, 400, 2, 1),
	)

	enemy.AddWeapon(weapons.NewBigSword(manager, enemy, "right"))

	if manager != nil {
		manager.AddObject(enemy)
	}
	return enemy
}

// ============================================================
// ОБНОВЛЕНИЕ
// ============================================================

func (enemy *OneBerserk) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {
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

func (b *OneBerserk) Draw(screen *ebiten.Image, camera *kamera.Camera) {

	b.StDraw(screen, camera)

}

// ============================================================
// ИНТЕРФЕЙС game.Entity
// ============================================================

func (b *OneBerserk) Tag() string {

	return b.TagName
}

func (b *OneBerserk) IsActive() bool {
	return true
}
