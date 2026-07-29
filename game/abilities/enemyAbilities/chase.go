package enemyabilities

import (
	"great-sword/game"
	"great-sword/game/hitboxes"
	"math"
)

var _ EnemyAbility = (*ChaseAbility)(nil)

type ChaseAbility struct {
	// Базовая скорость (когда враг далеко)
	BaseSpeed float64

	// Максимальная скорость (когда враг близко)
	MaxSpeed float64

	// Максимальная дистанция: начиная с этой дистанции враг ускоряется
	MaxDistance float64

	// Минимальная дистанция: на этой дистанции ускорение достигает максимума
	MinDistance float64

	// ТЕКУЩАЯ СКОРОСТЬ (рассчитывается каждый кадр)
	CurrentSpeed float64

	Active bool
}

func NewChaseAbility(baseSpeed, maxSpeed, maxDistance, minDistance float64) *ChaseAbility {
	return &ChaseAbility{
		BaseSpeed:   baseSpeed,
		MaxSpeed:    maxSpeed,
		MaxDistance: maxDistance,
		MinDistance: minDistance,
	}
}

func (c *ChaseAbility) Name() string {
	return "Chase"
}

// Update - рассчитывает скорость врага в зависимости от расстояния до игрока
// НЕ ДВИГАЕТ ВРАГА, только даёт ускорение (ДОБАВЛЯЕТ к текущей скорости)
func (c *ChaseAbility) Update(enemy EnemyUser, dt float64, manager *hitboxes.CollisionManager) bool {
	// Получаем позицию игрока (цель)
	tx, ty := enemy.GetTarget()
	ex, ey := enemy.GetPosition()

	// Вычисляем расстояние до игрока
	dx, dy := tx-ex, ty-ey
	dist := math.Sqrt(dx*dx + dy*dy)

	// ===== ПРОВЕРКА: ВРАГ В ЗОНЕ ДЕЙСТВИЯ =====
	// Если дистанция больше MaxDistance или меньше MinDistance → способность НЕ РАБОТАЕТ
	if dist > c.MaxDistance || dist < c.MinDistance {
		c.Active = false
		return false // способность не активна
	}

	// ===== РАСЧЁТ ДОБАВЛЯЕМОЙ СКОРОСТИ =====
	// dist = MaxDistance → addSpeed = 0 (базовая скорость)
	// dist = MinDistance → addSpeed = MaxSpeed - BaseSpeed (максимальное ускорение)
	// Чем ближе к MinDistance, тем больше ускорение
	normalizedDist := (dist - c.MinDistance) / (c.MaxDistance - c.MinDistance)
	// normalizedDist = 1 когда dist = MaxDistance (далеко)
	// normalizedDist = 0 когда dist = MinDistance (близко)

	// Добавляемая скорость (от 0 до MaxSpeed - BaseSpeed)
	addSpeed := (1 - normalizedDist) * (c.MaxSpeed - c.BaseSpeed)

	// ПРИМЕНЯЕМ: ДОБАВЛЯЕМ скорость к текущей скорости врага (не заменяем!)
	currentEnemySpeed := enemy.GetSpeed()
	enemy.SetSpeed(currentEnemySpeed + addSpeed)

	c.Active = true
	return true
}

func (c *ChaseAbility) Activate(enemy EnemyUser, world game.WorldView) bool {
	return true
}

func (c *ChaseAbility) IsActive() bool {
	return c.Active
}

func (c *ChaseAbility) GetCooldown() float64 {
	return 0
}

func (c *ChaseAbility) SetCooldown(value float64) {
	// Нет кулдауна
}
