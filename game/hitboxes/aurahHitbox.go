package hitboxes

import "math"

// CalculateAuraPush - рассчитывает силу отталкивания от ауры для двух объектов
// Параметры:
//   - a, b: объекты с интерфейсом HitBoxer
//
// Возвращает:
//   - pushX1, pushY1: сила отталкивания для объекта a
//   - pushX2, pushY2: сила отталкивания для объекта b
// В hitboxes/collision_manager.go

func CalculateAuraPush(a, b HitBoxer) (float64, float64, float64, float64) {
	// ===== ПРОВЕРКА: ЕСТЬ ЛИ АУРА У ОБЪЕКТОВ =====
	if !a.HasAura() && !b.HasAura() {
		return 0, 0, 0, 0
	}

	// ===== ПРОВЕРКА: РЕАГИРУЮТ ЛИ ОБЪЕКТЫ НА АУРУ =====
	if !a.AffectedByAura() && !b.AffectedByAura() {
		return 0, 0, 0, 0
	}

	// Получаем центры
	ax, ay, _, _ := a.GetAABB()
	bx, by, _, _ := b.GetAABB()

	dx := bx - ax
	dy := by - ay
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 0.01 {
		return 0, 0, 0, 0
	}

	// Нормализованное направление
	dirX := dx / dist
	dirY := dy / dist

	// Радиусы объектов
	_, _, ahw, ahh := a.GetAABB()
	_, _, bhw, bhh := b.GetAABB()
	radiusA := math.Max(ahw, ahh)
	radiusB := math.Max(bhw, bhh)

	// Расстояние между поверхностями
	surfaceDist := dist - (radiusA + radiusB)
	if surfaceDist < 0 {
		surfaceDist = 0
	}

	if surfaceDist > auraRadius {
		return 0, 0, 0, 0
	}

	// ===== РАСЧЁТ СИЛЫ С УЧЁТОМ ВЕСА =====
	weightA := a.GetWeight()
	weightB := b.GetWeight()

	// Если объект с большим весом сталкивается с лёгким
	// Лёгкий отталкивается, тяжёлый не чувствует ауру
	var pushForceA, pushForceB float64

	forceMultiplier := 1.0 - (surfaceDist / auraRadius)
	if forceMultiplier < 0 {
		forceMultiplier = 0
	}

	baseForce := auraForce * forceMultiplier

	// Если веса сильно отличаются (> 5 раз)
	if weightA > weightB*5 {
		// Тяжёлый A не чувствует ауру от лёгкого B
		pushForceA = 0
		pushForceB = baseForce * (weightA / weightB) // лёгкий отлетает
	} else if weightB > weightA*5 {
		pushForceA = baseForce * (weightB / weightA) // лёгкий отлетает
		pushForceB = 0
	} else {
		// Оба объекта чувствуют ауру
		totalWeight := weightA + weightB
		pushForceA = baseForce * (weightB / totalWeight)
		pushForceB = baseForce * (weightA / totalWeight)
	}

	// Применяем только если объект реагирует на ауру
	if !a.AffectedByAura() {
		pushForceA = 0
	}
	if !b.AffectedByAura() {
		pushForceB = 0
	}

	return -dirX * pushForceA, -dirY * pushForceA,
		dirX * pushForceB, dirY * pushForceB
}
