package weapons

import (
	"great-sword/game"
	effectsmass "great-sword/game/effects/effectsMass"
	"great-sword/game/hitboxes"
	"math/rand"
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
			90,  // ширина
			40,  // высота
			2.5, // плотность (оригинальная)
			3,   // проризание из 10
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
				effectsmass.NewDamageEffect(5),
			},
			reflect.TypeOf((*game.PlayerLegInter)(nil)).Elem(),
		),
	}

	manager.AddObject(weapon)

	return weapon
}

func (b *LitleKnife) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {

	b.StandartUpdate()
	return false
}

func (b *LitleKnife) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	b.StandartDraw(screen, camera)
}

func (b *LitleKnife) Tag() string {
	return b.ID
}
