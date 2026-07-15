package abilities

import (
	"great-sword/game"
	"great-sword/game/common"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Ability = (*Dash)(nil)

var lastMoveX, lastMoveY float64

type Dash struct {
	Active      bool
	Timer       float64
	Duration    float64
	Cooldown    float64
	CooldownMax float64
	Speed       float64
}

func NewDash() *Dash {
	return &Dash{
		Active:      false,
		Duration:    0.3,
		CooldownMax: 2.0,
		Cooldown:    2.0,
		Speed:       1200,
	}
}

func (d *Dash) Name() string {
	return "Dash"
}

func (d *Dash) Update(world game.WorldView) bool {
	dt := 1.0 / 60.0

	// Уменьшаем кулдаун
	if d.Cooldown > 0 {
		d.Cooldown -= dt
	}

	// Если рывок активен
	if d.Active {
		d.Timer -= dt
		for _, pleg := range world.SearchEntities("playerLeg") {
			if p, ok := pleg.(game.PlayerLegInter); ok {

				//устанавливаем максимальную скорость
				p.SetMaxSpeed(d.Speed * 1.5)

				moveX, moveY := p.GetDirection()

				// Применяем рывок (ускорение)
				p.ApplyForce(moveX*d.Speed, moveY*d.Speed)

				if d.Timer <= 0 {
					d.Active = false
					// Восстанавливаем скорость
					p.SetSpeed(lastMoveX, lastMoveY)
				}
			}
		}
		return true
	}
	return false
}

func (d *Dash) Activate(world game.WorldView) {
	// Проверяем кулдаун
	if d.Cooldown > 0 {
		return
	}

	// Ищем игрока через интерфейс
	for _, pleg := range world.SearchEntities("playerLeg") {
		p, ok := pleg.(game.PlayerLegInter)
		if !ok {
			continue
		}

		// Получаем направление движения
		moveX, moveY := p.GetDirection()

		lastMoveX, lastMoveY = p.GetSpeed()

		// Если не двигается — рывок по углу
		if moveX == 0 && moveY == 0 {
			moveX = float64(math.Cos(0))
			moveY = float64(math.Sin(0))
		}

	}

	// Активируем способность
	d.Active = true
	d.Timer = d.Duration
	if !common.SwordExist {
		d.Cooldown = d.CooldownMax / 2
	} else {
		d.Cooldown = d.CooldownMax
	}
}

func (d *Dash) IsActive() bool {
	return d.Active
}

func (d *Dash) CooldownAbilites() float64 {
	return d.Cooldown
}

func (d *Dash) SetCooldown(value float64) {
	d.Cooldown = value
}

func (d *Dash) Draw(screen *ebiten.Image) {
	// Можно отображать индикатор кулдауна
}
