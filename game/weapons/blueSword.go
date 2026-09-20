package weapons

import (
	"great-sword/game"
	"great-sword/game/common"
	effectsmass "great-sword/game/effects/effectsMass"
	"great-sword/game/hitboxes"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

var _ Weapon = (*BlueSwordWeapon)(nil)
var _ hitboxes.HitBoxer = (*BlueSwordWeapon)(nil)
var _ hitboxes.LetterSender = (*BlueSwordWeapon)(nil)
var _ hitboxes.RotatableHitBoxer = (*BlueSwordWeapon)(nil)

// ============================================================
// BlueSwordWeapon - копия оригинального BlueSword на BaseWeapon
// ============================================================

type BlueSwordWeapon struct {
	*BaseWeapon
}

// NewBlueSwordWeapon - создаёт копию BlueSword на основе BaseWeapon
func NewBlueSwordWeapon(manager *hitboxes.CollisionManager, user WeaponUser) *BlueSwordWeapon {
	// Оригинальные параметры BlueSword:
	// - Weight: 7
	// - Density: 5.5
	// - Размер: common.SwordAttachmentWidth x common.SwordAttachmentHeight
	// - Письма: BurnEffect (2.0 сек, 50 урона, 3 передачи)
	// - Цель: game.Enemy

	weapon := &BlueSwordWeapon{
		BaseWeapon: NewBaseWeapon(
			common.SwordAttachmentWidth,  // ширина
			common.SwordAttachmentHeight, // высота
			10.5,                         // плотность (оригинальная)
			7,                            // проризание из 10
			WeaponPhysics{
				Weight:              7.0,   // оригинальный вес
				Friction:            0.5,   // стандартное трение
				ImpactForce:         15.0,  // сила удара
				MaxAngularVelocity:  360.0, // макс скорость вращения
				AngularAcceleration: 720.0, // ускорение вращения
				Damping:             0.9,   // демпфирование
				SmoothingPos:        0.1,   //сила пиритаскивания оружия
				SmoothingReturn:     0.1,   // скорость возврата угла к цели
			},
			"blueSwordCopy",
			user,
		),
	}

	weapon.Letters = []*hitboxes.Letter{
		hitboxes.NewLetter(
			true,
			0.5,
			[]hitboxes.Effect{
				effectsmass.NewDamageEffect(50),
			},
			reflect.TypeOf((*game.Enemy)(nil)).Elem(),
		),
	}

	manager.AddObject(weapon)

	return weapon
}

func (b *BlueSwordWeapon) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {

	b.StandartUpdate()

	return false
}

func (b *BlueSwordWeapon) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	b.StandartDraw(screen, camera)
}

func (b *BlueSwordWeapon) Tag() string {
	return b.ID
}
