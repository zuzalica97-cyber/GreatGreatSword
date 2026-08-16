package effectsmass

import (
	"great-sword/game/hitboxes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

// ============================================================
// ЭФФЕКТ НАНЕСЕНИЯ УРОНА (живёт 1 кадр)
// ============================================================

type DamageEffect struct {
	BaseEffect
	damage float64
}

// NewDamageEffect - создаёт эффект урона
//   - damage: количество урона
func NewDamageEffect(damage float64) *DamageEffect {
	return &DamageEffect{
		BaseEffect: NewBaseEffect(
			"damage", // тип
			0,        // duration = 0 (живёт 1 кадр)
			0,        // maxTransfers = 0 (не передаётся)
			false,    // canStack = false (не суммируется)
			false,    // canExtend = false (не продлевается)
		),
		damage: damage,
	}
}

// Apply - применяет урон к цели
func (d *DamageEffect) Apply(target hitboxes.EffectUser) {
}

// Update - обновляет эффект (всегда завершается за 1 кадр)
// Возвращает true (эффект завершён)
func (d *DamageEffect) Update(target hitboxes.EffectUser, dt float64) bool {
	// Наносим урон при первом обновлении
	target.TakeDamage(d.damage)
	// Завершаем эффект
	return true
}

// OnTransfer - вызывается при передаче эффекта
func (d *DamageEffect) OnTransfer(newTarget hitboxes.EffectUser) {
	// Не передаётся
}

func (b *DamageEffect) Draw(screen *ebiten.Image, camera *kamera.Camera, target hitboxes.EffectUser) {
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
		color.RGBA{255, 0, 0, 100},
		true,
	)

}

// Clone - создаёт копию эффекта
func (d *DamageEffect) Clone() hitboxes.Effect {
	return &DamageEffect{
		BaseEffect: d.BaseEffect,
		damage:     d.damage,
	}
}
