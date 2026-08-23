// effects/burn.go

package effectsmass

import (
	"great-sword/game/hitboxes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

type BurnEffect struct {
	BaseEffect
	damagePerSecond float64
	minTimer        float64
}

func NewBurnEffect(duration, damagePerSecond float64, maxTransfers int) *BurnEffect {
	return &BurnEffect{
		BaseEffect: NewBaseEffect(
			"burn",       // тип
			duration,     // длительность
			maxTransfers, // передачи
			false,        // canStack (нельзя больше одного горения)
			true,         // canExtend (можно продлить)
		),
		damagePerSecond: damagePerSecond,
	}
}

func (b *BurnEffect) Apply(target hitboxes.EffectUser) {
	target.SetSpeed(target.GetSpeed() * 0.8)
}

func (b *BurnEffect) OnTransfer(newTarget hitboxes.EffectUser) {
	b.remainingTransfers--
}

func (b *BurnEffect) Update(target hitboxes.EffectUser, dt float64) bool {
	if b.UpdateBase(dt) {
		return true // эффект завершён
	}

	if b.minTimer <= 0 {
		b.minTimer = 0.1
		target.TakeDamage(b.damagePerSecond * 0.1) //Дз Доделать нормально чтобы работало а ппотом комитт!!!
		target.SetSpeed(target.GetSpeed() * 0.8)
	} else {
		b.minTimer -= dt
	}

	return false
}

func (b *BurnEffect) Draw(screen *ebiten.Image, camera *kamera.Camera, target hitboxes.EffectUser) {

	enemyX, enemyY := target.GetPosition()

	screenX := enemyX - camera.X
	screenY := enemyY - camera.Y

	size := target.GetSize()

	vector.FillRect(
		screen,
		float32(screenX),
		float32(screenY),
		float32(size),
		float32(size),
		color.RGBA{230, 20, 0, 100},
		true,
	)

}

func (b *BurnEffect) Clone() hitboxes.Effect {
	return &BurnEffect{
		BaseEffect: BaseEffect{
			id:                 b.id + "_copy", // новый ID
			effectType:         b.effectType,
			active:             true,
			duration:           b.duration, // ← оригинальная длительность!
			timer:              b.duration, // ← СБРАСЫВАЕМ ТАЙМЕР!
			maxTransfers:       b.maxTransfers,
			remainingTransfers: b.remainingTransfers,
			canStack:           b.canStack,
			canExtend:          b.canExtend,
		},
		damagePerSecond: b.damagePerSecond,
	}
}
