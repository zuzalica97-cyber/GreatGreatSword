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

var _ game.Entity = (*Berseks)(nil)
var _ hitboxes.HitBoxer = (*OneBerserk)(nil)
var _ enemyabilities.EnemyUser = (*OneBerserk)(nil)

// ============================================================
// ОСНОВНАЯ СТРУКТУРА ВРАГА
// ============================================================

type Berseks struct {
	BerserkMass []*OneBerserk
}

type OneBerserk struct {
	*BaseEnemy
	Abilities []enemyabilities.EnemyAbility
	oldSpeed  float64
}

// ============================================================
// КОНСТРУКТОР
// ============================================================

func NewBerserk() *Berseks {
	return &Berseks{
		BerserkMass: make([]*OneBerserk, 0),
	}
}

// ============================================================
// СОЗДАНИЕ ВРАГА
// ============================================================

func (b *Berseks) Spawn(x, y float64, manager *hitboxes.CollisionManager) {
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

	enemy.SetWeapon(weapons.NewBigSword(manager, enemy))

	b.BerserkMass = append(b.BerserkMass, enemy)
	if manager != nil {
		manager.AddObject(enemy)
	}
}

// ============================================================
// ОБНОВЛЕНИЕ
// ============================================================

func (b *Berseks) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {
	dt := 1.0 / 60.0

	playerX, playerY := getPlayerPosition(worldView)

	if len(b.BerserkMass) < 1 {
		x, y := RangomSpawnInWall(50)
		b.Spawn(x, y, manager)
	}

	for i := 0; i < len(b.BerserkMass); i++ {
		enemy := b.BerserkMass[i]

		// Проверка смерти
		var died bool
		b.BerserkMass, died = DeathScan(manager, enemy.BaseEnemy, b.BerserkMass, i)
		if died {
			i-- // корректируем индекс после удаления
			continue
		}

		// === ОБНОВЛЕНИЕ КУЛДАУНА ===
		enemy.StUpdateCoolDown(dt, worldView, manager)

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

	}

	return false
}

// ============================================================
// ОТРИСОВКА
// ============================================================

func (b *Berseks) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	for _, enemy := range b.BerserkMass {
		enemy.StDraw(screen, camera)
	}
}

// ============================================================
// ИНТЕРФЕЙС game.Entity
// ============================================================

func (b *Berseks) Tag() string {
	return "berserk"
}

func (b *Berseks) IsActive() bool {
	return true
}
