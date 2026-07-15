package abilities

import (
	"great-sword/game"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Ability = (*FastRotate)(nil)

type FastRotate struct {
	Active        bool
	Timer         float64
	Duration      float64
	Cooldown      float64
	CooldownMax   float64
	Speed         float64
	OriginalSpeed float64
	OriginalAccel float64
	OriginalDecel float64
}

func NewFastRotate() *FastRotate {
	return &FastRotate{
		Active:      false,
		Duration:    0.5,
		CooldownMax: 2.0,
		Cooldown:    2.0,
		Speed:       2000,
	}
}

func (f *FastRotate) Name() string {
	return "FastRotate"
}

func (f *FastRotate) Update(world game.WorldView) bool {
	dt := 1.0 / 60.0

	if f.Cooldown > 0 {
		f.Cooldown -= dt
	}

	if f.Active {

		f.Timer -= dt

		// ПРИМЕНЯЕМ эффект (устанавливаем параметры)
		for _, entity := range world.SearchEntities("playerHead") {
			if head, ok := entity.(game.PlayerHeadInter); ok {
				// Устанавливаем быстрые параметры
				head.SetCurrentRotSpeed(f.Speed)
				head.SetDeAceleration(0) // не тормозит
				axel := head.GetAceleration()

				if axel < 500 {
					axel = 500
				}

				head.SetAceleration(axel * 2) // быстро разгоняется
			}
		}

		if f.Timer <= 0 {
			f.Active = false
			// ВОССТАНАВЛИВАЕМ параметры
			for _, entity := range world.SearchEntities("playerHead") {
				if head, ok := entity.(game.PlayerHeadInter); ok {
					head.SetCurrentRotSpeed(f.OriginalSpeed)
					head.SetAceleration(f.OriginalAccel)
					head.SetDeAceleration(f.OriginalDecel)
				}
			}
		}
		return true
	}
	return false
}

func (f *FastRotate) Activate(world game.WorldView) {
	if f.Cooldown > 0 {
		return
	}

	// Сохраняем оригинальные параметры
	for _, entity := range world.SearchEntities("playerHead") {
		if head, ok := entity.(game.PlayerHeadInter); ok {
			f.OriginalSpeed = head.GetCurrentRotSpeed()
			f.OriginalAccel = head.GetAceleration()
			f.OriginalDecel = head.GetDeAceleration()
		}
	}

	f.Active = true
	f.Timer = f.Duration
	f.Cooldown = f.CooldownMax
}

func (f *FastRotate) IsActive() bool {
	return f.Active
}

func (f *FastRotate) CooldownAbilites() float64 {
	return f.Cooldown
}

func (f *FastRotate) SetCooldown(value float64) {
	f.Cooldown = value
}

func (f *FastRotate) Draw(screen *ebiten.Image) {
}
