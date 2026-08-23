package enemy

import (
	"great-sword/game"
	enemyabilities "great-sword/game/abilities/enemyAbilities"
	"great-sword/game/common"
	effectsmass "great-sword/game/effects/effectsMass"
	"great-sword/game/hitboxes"
	"great-sword/game/player"
	"image/color"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
			200,                         // health
			10,                          // damage
			200,                         // baseSpeed
			400,                         // maxSpeed
			color.RGBA{90, 90, 90, 255}, //ДЗ нужно поравить силу ооталкивания чтобы тяжёлые обьекты легко отталкивали лёгкие а то аура мешает
			50,
			1,
			"berserk",
		),
		oldSpeed: 0,
	}

	enemy.Abilities = append(enemy.Abilities,
		enemyabilities.NewDashAbility(800, 400, 2, 1),
	)

	enemy.Letters = []*hitboxes.Letter{
		hitboxes.NewLetter(
			true,
			0.5,
			[]hitboxes.Effect{
				effectsmass.NewDamageEffect(float64(enemy.Damage)),
			},
			reflect.TypeOf((*game.PlayerLegInter)(nil)).Elem(),
		),
	}

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

		// === ПРОВЕРКА СМЕРТИ ===
		if !enemy.IsActive() || enemy.GetHealth() <= 0 {
			if manager != nil {
				manager.RemoveObject(enemy)
			}
			b.BerserkMass[i] = nil
			b.BerserkMass = append(b.BerserkMass[:i], b.BerserkMass[i+1:]...)
			i--
			common.Score++
			player.ActivateBoost()
			continue
		}

		// === ОБНОВЛЕНИЕ КУЛДАУНА ===
		enemy.UpdateCooldown(dt)

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

		friction := 0.9999   // чем ближе к 1, тем дольше скользит
		acceleration := 50.0 // как быстро разгоняется

		newSpeed, newDirX, newDirY := enemy.EnemySlideMovmentFunc(friction, acceleration, dt)

		enemy.CurrentSpeed = newSpeed
		enemy.SetDirection(newDirX, newDirY)

		// === ДВИЖЕНИЕ ===
		newX, newY := MoveEnemyToTareget(enemy.BaseEnemy, dt)

		enemy.SetPosition(newX, newY)

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
		screenX := enemy.X - camera.X
		screenY := enemy.Y - camera.Y

		Color := enemy.Color
		if enemy.CooldownActive {
			Color = color.RGBA{90, 90, 90, 255}
		}
		for _, ability := range enemy.Abilities {
			if ability.Name() == "Dash" && ability.IsActive() {
				Color = color.RGBA{140, 90, 90, 225}
				break
			}
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

		for _, effect := range enemy.Effects {
			effect.Draw(screen, camera, enemy)
		}
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
