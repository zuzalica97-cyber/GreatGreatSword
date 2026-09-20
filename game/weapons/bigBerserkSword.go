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

type BigSword struct {
	*BaseWeapon
}

// NewBlueSwordWeapon - создаёт копию BlueSword на основе BaseWeapon
func NewBigSword(manager *hitboxes.CollisionManager, user WeaponUser) *BigSword {
	// Оригинальные параметры BlueSword:
	// - Weight: 7
	// - Density: 5.5
	// - Размер: common.SwordAttachmentWidth x common.SwordAttachmentHeight
	// - Письма: BurnEffect (2.0 сек, 50 урона, 3 передачи)
	// - Цель: game.Enemy

	weapon := &BigSword{
		BaseWeapon: NewBaseWeapon(
			250, // ширина
			80,  // высота
			4,   // плотность (оригинальная)
			4,   // проризание из 10
			WeaponPhysics{
				Weight:              10.0,  // оригинальный вес
				Friction:            0.5,   // стандартное трение
				ImpactForce:         10.0,  // сила удара
				MaxAngularVelocity:  200.0, // макс скорость вращения
				AngularAcceleration: 400.0, // ускорение вращения
				Damping:             0.4,   // демпфирование
				SmoothingPos:        0.5,   //сила пиритаскивания оружия
				SmoothingReturn:     0.05,  // скорость возврата угла к цели
			},
			"bigSword"+string(rand.Intn(1000)),
			user,
		),
	}

	weapon.Letters = []*hitboxes.Letter{
		hitboxes.NewLetter(
			true,
			0.5,
			[]hitboxes.Effect{
				effectsmass.NewDamageEffect(20),
			},
			reflect.TypeOf((*game.PlayerLegInter)(nil)).Elem(),
		),
	}

	manager.AddObject(weapon)

	return weapon
}

func (b *BigSword) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {

	b.StandartUpdate()
	return false
}

func (b *BigSword) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	b.StandartDraw(screen, camera)
}

func (b *BigSword) Tag() string {
	return b.ID
}
