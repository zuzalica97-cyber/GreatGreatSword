package enemyabilities

import (
	"great-sword/game"
	"great-sword/game/hitboxes"
	"math"
)

var _ EnemyAbility = (*FleeAbility)(nil)

type FleeAbility struct {
	// Безопасная дистанция (когда враг достиг её, он останавливается)
	SafeDistance float64

	// Максимальная скорость убегания (когда игрок очень близко)
	MaxSpeed float64

	// Максимальная дистанция: начиная с этой дистанции враг убегает
	MaxDistance float64

	// Минимальная дистанция: на этой дистанции скорость убегания достигает максимума
	MinDistance float64

	// ТЕКУЩАЯ СКОРОСТЬ (рассчитывается каждый кадр)
	CurrentSpeed float64

	Active bool
}

func NewFleeAbility(safeDistance, maxSpeed, maxDistance, minDistance float64) *FleeAbility { //ДЗ нужно сделать плавное передвижение врагои колизию пуль. потом нужно сделать отдельную структуру пуль , потом приступить к эффектам
	return &FleeAbility{
		SafeDistance: safeDistance,
		MaxSpeed:     maxSpeed,
		MaxDistance:  maxDistance,
		MinDistance:  minDistance,
	}
}

func (f *FleeAbility) Name() string {
	return "Flee"
}

// Update - рассчитывает скорость убегания
// На SafeDistance скорость = 0 (враг стоит)
// Чем ближе к MinDistance, тем быстрее убегает
// Чем дальше от SafeDistance, тем быстрее возвращается к SafeDistance
func (f *FleeAbility) Update(enemy EnemyUser, dt float64, manager *hitboxes.CollisionManager) bool {
	// Получаем позицию игрока
	tx, ty := enemy.GetTarget()
	ex, ey := enemy.GetPosition()

	// Вычисляем расстояние до игрока
	dx, dy := tx-ex, ty-ey
	dist := math.Sqrt(dx*dx + dy*dy)

	// Если дистанция больше MaxDistance → способность НЕ РАБОТАЕТ
	if dist > f.MaxDistance {
		f.Active = false
		return false
	}

	// Если дистанция равна SafeDistance → скорость 0 (стоим на месте)
	if math.Abs(dist-f.SafeDistance) < 0.1 {
		f.CurrentSpeed = 0
		f.Active = true
		enemy.SetSpeed(f.CurrentSpeed)
		return true
	}

	// Если дистанция меньше SafeDistance → убегаем (отрицательная скорость)
	// Чем ближе к MinDistance, тем быстрее убегаем
	if dist < f.SafeDistance {
		// dist = SafeDistance → speed = 0
		// dist = MinDistance → speed = -MaxSpeed
		normalizedDist := (dist - f.MinDistance) / (f.SafeDistance - f.MinDistance)
		if normalizedDist < 0 {
			normalizedDist = 0
		}
		if normalizedDist > 1 {
			normalizedDist = 1
		}
		// От 0 до -MaxSpeed
		f.CurrentSpeed = -f.MaxSpeed * (1 - normalizedDist)
		f.Active = true
		enemy.SetSpeed(f.CurrentSpeed)
		return true
	}

	// Если дистанция больше SafeDistance → возвращаемся к SafeDistance (положительная скорость)
	// dist = SafeDistance → speed = 0
	// dist = MaxDistance → speed = MaxSpeed (возвращается)
	if dist > f.SafeDistance {
		normalizedDist := (dist - f.SafeDistance) / (f.MaxDistance - f.SafeDistance)
		if normalizedDist < 0 {
			normalizedDist = 0
		}
		if normalizedDist > 1 {
			normalizedDist = 1
		}
		f.CurrentSpeed = f.MaxSpeed * normalizedDist
		f.Active = true
		enemy.SetSpeed(f.CurrentSpeed)
		return true
	}

	f.Active = false
	return false
}

func (f *FleeAbility) Activate(enemy EnemyUser, world game.WorldView) bool {
	return true
}

func (f *FleeAbility) IsActive() bool {
	return f.Active
}

func (f *FleeAbility) GetCooldown() float64 {
	return 0
}

func (f *FleeAbility) SetCooldown(value float64) {
	// Нет кулдауна
}
