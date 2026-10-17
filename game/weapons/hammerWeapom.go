package weapons

import (
	"great-sword/game"
	effectsmass "great-sword/game/effects/effectsMass"
	"great-sword/game/hitboxes"
	"image/color"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

var _ Weapon = (*Hammer)(nil)
var _ hitboxes.HitBoxer = (*Hammer)(nil)
var _ hitboxes.LetterSender = (*Hammer)(nil)
var _ hitboxes.RotatableHitBoxer = (*Hammer)(nil)

// ============================================================
// HAMMER - молоток, состоящий из рукоятки и навершия
// ============================================================

type Hammer struct {
	*BaseWeapon
}

func NewHammer(user WeaponUser, side string, manager *hitboxes.CollisionManager) *Hammer { //ДЗ исправиь привязку чтобы части  не отваливадись и с исправь отправку писем
	// Создаём базовое оружие
	// Ширина и высота — это общий bounding box (не используется для коллизий)
	// Реальные коллизии — у частей (Parts)
	weapon := &Hammer{
		BaseWeapon: NewBaseWeapon(
			30,  // ширина рукоятки
			180, // высота (общая)
			5.0, // плотность
			0.8, // проризание
			WeaponPhysics{
				Weight:              10.0,
				Friction:            0.5,
				ImpactForce:         25.0,
				MaxAngularVelocity:  360.0,
				AngularAcceleration: 720.0,
				Damping:             0.9,
				SmoothingReturn:     0.15,
				SmoothingPos:        0.2,
			},
			side, // сторона привязки
			"hammer",
			user,
		),
	}
	// ===== НАВЕРШИЕ =====
	head := NewWeaponPart(
		0,    // смещение X (по центру)
		-110, // смещение Y (вверх от центра)
		80,   // ширина
		40,   // высота
		0,    // угол
		"hammer_head",
		color.RGBA{100, 100, 100, 255},
		weapon.BaseWeapon,
	)
	head.Weight = 7.0
	head.Density = 1.0 // тяжёлое, твёрдое

	// ===== ПИСЬМА ДЛЯ НАВЕРШИЯ =====
	head.Letters = []*hitboxes.Letter{
		hitboxes.NewLetter(
			true, // бесконечное
			0.5,  // кулдаун
			[]hitboxes.Effect{
				effectsmass.NewDamageEffect(20),
			},
			reflect.TypeOf((*game.PlayerLegInter)(nil)).Elem(),
		),
	}

	// Добавляем части в оружие
	weapon.Parts = []*WeaponPart{head}

	// Привязываем к пользователю
	weapon.Attach(user)

	manager.AddObject(weapon)
	for _, part := range weapon.Parts {
		manager.AddObject(part)
	}

	return weapon
}

func (b *Hammer) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {

	b.StandartUpdate()

	return false
}

func (b *Hammer) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	b.StandartDraw(screen, camera)
}

func (b *Hammer) Tag() string {
	return b.ID
}
