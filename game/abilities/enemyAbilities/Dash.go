package enemyabilities

import (
	"great-sword/game"
	"great-sword/game/hitboxes"
	"math"
)

var _ EnemyAbility = (*DashAbility)(nil)

type DashAbility struct {
	Active        bool
	Timer         float64
	Duration      float64
	Cooldown      float64
	CooldownMax   float64
	Speed         float64
	DashDistance  float64 // дистанция активации
	DirX, DirY    float64 // фиксированное направление
	OriginalSpeed float64 // для восстановления
}

func NewDashAbility(dashDistance, speed, cooldownMax, diuration float64) *DashAbility {
	return &DashAbility{
		Active:       false,
		Duration:     diuration,
		CooldownMax:  cooldownMax,
		Speed:        speed,
		DashDistance: dashDistance,
	}
}

func (d *DashAbility) Name() string { return "Dash" }

func (d *DashAbility) Update(enemy EnemyUser, dt float64, manager *hitboxes.CollisionManager) bool {
	if d.Cooldown > 0 {
		d.Cooldown -= dt
	}

	if d.Active {
		d.Timer -= dt

		enemy.SetDirection(d.DirX, d.DirY)
		enemy.SetSpeed(d.Speed)

		if d.Timer <= 0 {
			d.Active = false
		}
		return true
	}
	return false
}

func (d *DashAbility) Activate(enemy EnemyUser, world game.WorldView) bool {
	if d.Cooldown > 0 {
		return false
	}

	// Получаем позицию игрока (цель)
	tx, ty := enemy.GetTarget()
	ex, ey := enemy.GetPosition()

	// Вычисляем направление к игроку
	dx, dy := tx-ex, ty-ey
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > 0.01 && dist < d.DashDistance {
		// ФИКСИРУЕМ направление в момент активации
		d.DirX = dx / dist
		d.DirY = dy / dist

		// Активируем рывок
		d.Active = true
		d.Timer = d.Duration
		d.Cooldown = d.CooldownMax

		// Устанавливаем скорость рывка
		enemy.SetSpeed(d.Speed)
		return true
	}
	return false
}

func (d *DashAbility) IsActive() bool { return d.Active }

func (d *DashAbility) GetCooldown() float64 { return d.Cooldown }

func (d *DashAbility) SetCooldown(value float64) { d.Cooldown = value }

func (d *DashAbility) SetDirection(x, y float64) {
	d.DirX = x
	d.DirY = y
}
