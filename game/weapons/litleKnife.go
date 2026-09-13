package weapons

import (
	"great-sword/game"
	effectsmass "great-sword/game/effects/effectsMass"
	"great-sword/game/hitboxes"
	"image/color"
	"math"
	"math/rand"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

var _ Weapon = (*BlueSwordWeapon)(nil)
var _ hitboxes.HitBoxer = (*BlueSwordWeapon)(nil)
var _ hitboxes.LetterSender = (*BlueSwordWeapon)(nil)
var _ hitboxes.RotatableHitBoxer = (*BlueSwordWeapon)(nil)

// ============================================================
// BlueSwordWeapon - копия оригинального BlueSword на BaseWeapon
// ============================================================

type LitleKnife struct {
	*BaseWeapon
}

// NewBlueSwordWeapon - создаёт копию BlueSword на основе BaseWeapon
func NewlitleKnife(manager *hitboxes.CollisionManager, user WeaponUser) *LitleKnife {
	// Оригинальные параметры BlueSword:
	// - Weight: 7
	// - Density: 5.5
	// - Размер: common.SwordAttachmentWidth x common.SwordAttachmentHeight
	// - Письма: BurnEffect (2.0 сек, 50 урона, 3 передачи)
	// - Цель: game.Enemy

	weapon := &LitleKnife{
		BaseWeapon: NewBaseWeapon(
			90,   // ширина
			40,   // высота
			2.5,  // плотность (оригинальная)
			0.03, // проризание
			WeaponPhysics{
				Weight:              3.0,   // оригинальный вес
				Friction:            0.5,   // стандартное трение
				ImpactForce:         5.0,   // сила удара
				MaxAngularVelocity:  280.0, // макс скорость вращения
				AngularAcceleration: 500.0, // ускорение вращения
				Damping:             0.9,   // демпфирование
				SmoothingPos:        0.5,   //сила пиритаскивания оружия
				SmoothingReturn:     0.1,   // скорость возврата угла к цели
			},
			"litleKnife"+string(rand.Intn(1000)),
			user,
		),
	}

	weapon.Letters = []*hitboxes.Letter{
		hitboxes.NewLetter(
			true,
			0.5,
			[]hitboxes.Effect{
				effectsmass.NewDamageEffect(0),
			},
			reflect.TypeOf((*game.PlayerLegInter)(nil)).Elem(),
		),
	}

	manager.AddObject(weapon)

	return weapon
}

func (b *LitleKnife) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {

	dt := 1.0 / 60.0

	for _, letter := range b.Letters {
		letter.UpdateCoolDown(dt)
	}

	// 1. Обновляем целевую позицию (для ориентации)
	b.UpdateAttachmentTarget()

	// 2. Применяем физику "гвоздя" с инерцией
	b.UpdateWeaponAngle(dt)

	b.UpdateWeaponPosition(dt)

	return false
}

func (b *LitleKnife) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	Color := color.RGBA{0, 100, 200, 255}

	// Создаём временное изображение для меча
	swordImg := ebiten.NewImage(int(b.Width), int(b.Height))
	swordImg.Fill(Color)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-b.Width/2, -b.Height/2)
	op.GeoM.Rotate(b.Angle * math.Pi / 180)
	op.GeoM.Translate(b.PositionX, b.PositionY) // мировые координаты

	// Камера сама применит смещение
	camera.Draw(swordImg, op, screen)

	// Отрисовка точки контакта
	if b.DebugContactActive {
		size := 16.0
		screenX := float32(b.DebugContactX - camera.X - size/2)
		screenY := float32(b.DebugContactY - camera.Y - size/2)

		vector.FillRect(
			screen,
			screenX, screenY,
			float32(size), float32(size),
			color.RGBA{255, 0, 0, 255}, // красный
			true,
		)
		b.DebugContactActive = false // сбрасываем после отрисовки
	}
}

func (b *LitleKnife) Tag() string {
	return b.ID
}
