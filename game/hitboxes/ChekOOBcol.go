package hitboxes

import (
	"math"

	"github.com/setanarut/coll"
	v "github.com/setanarut/v"
)

// / checkOBBvsOBB - проверка столкновения двух вращающихся объектов (SAT)
func (cm *CollisionManager) checkOBBvsOBB(rot1, rot2 RotatableHitBoxer) {
	cx1, cy1, hw1, hh1, angle1 := rot1.GetOBB()
	cx2, cy2, hw2, hh2, angle2 := rot2.GetOBB()

	overlap, normal, penetration := SAT_OBBvsOBB(cx1, cy1, hw1, hh1, angle1, cx2, cy2, hw2, hh2, angle2)

	if !overlap {
		return
	}

	hit := &coll.Hit{
		Normal: normal,
		Data:   penetration,
	}

	contactX := (cx1 + cx2) / 2
	contactY := (cy1 + cy2) / 2

	rot1.Debag(contactX, contactY)
	rot2.Debag(contactX, contactY)

	cm.resolveOBBvsOBB(rot1, rot2, hit)
}

func SAT_OBBvsOBB(
	cx1, cy1, hw1, hh1, angle1 float64,
	cx2, cy2, hw2, hh2, angle2 float64,
) (bool, v.Vec, float64) {

	rad1 := angle1 * math.Pi / 180
	rad2 := angle2 * math.Pi / 180

	// 4 оси для проверки
	axes := []v.Vec{
		{X: math.Cos(rad1), Y: math.Sin(rad1)},  // ось X OBB1
		{X: -math.Sin(rad1), Y: math.Cos(rad1)}, // ось Y OBB1
		{X: math.Cos(rad2), Y: math.Sin(rad2)},  // ось X OBB2
		{X: -math.Sin(rad2), Y: math.Cos(rad2)}, // ось Y OBB2
	}

	dx := cx2 - cx1
	dy := cy2 - cy1

	minOverlap := math.MaxFloat64
	var minAxis v.Vec

	for _, axis := range axes {
		// Проекция расстояния между центрами
		distProj := math.Abs(dx*axis.X + dy*axis.Y)

		// Проекция первого OBB
		proj1 := hw1*math.Abs(axis.X*math.Cos(rad1)+axis.Y*math.Sin(rad1)) +
			hh1*math.Abs(-axis.X*math.Sin(rad1)+axis.Y*math.Cos(rad1))

		// Проекция второго OBB
		proj2 := hw2*math.Abs(axis.X*math.Cos(rad2)+axis.Y*math.Sin(rad2)) +
			hh2*math.Abs(-axis.X*math.Sin(rad2)+axis.Y*math.Cos(rad2))

		// Перекрытие
		overlap := proj1 + proj2 - distProj

		// Если нет перекрытия — нет столкновения
		if overlap < 0 {
			return false, v.Vec{}, 0
		}

		// Ищем минимальное перекрытие
		if overlap < minOverlap {
			minOverlap = overlap
			minAxis = axis
		}
	}

	// Нормаль — от первого ко второму
	if dx*minAxis.X+dy*minAxis.Y < 0 {
		minAxis.X = -minAxis.X
		minAxis.Y = -minAxis.Y
	}

	return true, minAxis, minOverlap
}

func (cm *CollisionManager) resolveOBBvsOBB(rot1, rot2 RotatableHitBoxer, hit *coll.Hit) {

	// ===== ОБРАБОТКА СТОЛКНОВЕНИЯ =====
	normal := hit.Normal
	penetration := hit.Data

	// ===== 1. ОТТАЛКИВАНИЕ ОТ СТОЛКНОВЕНИЯ =====
	weight1 := rot1.GetWeight()
	weight2 := rot2.GetWeight()
	totalWeight := weight1 + weight2
	if totalWeight == 0 {
		totalWeight = 1
	}

	// Сила отталкивания (базовая)
	pushStrength := 0.8

	// Коэффициенты отталкивания (чем больше вес, тем меньше отталкивание)
	pushFactor1 := weight2 / totalWeight
	pushFactor2 := weight1 / totalWeight

	if !rot1.IsStatic() {
		pushX := -normal.X * penetration * pushStrength * pushFactor1
		pushY := -normal.Y * penetration * pushStrength * pushFactor1
		rot1.ApplyPush(pushX, pushY)
	}
	if !rot2.IsStatic() {
		pushX := normal.X * penetration * pushStrength * pushFactor2
		pushY := normal.Y * penetration * pushStrength * pushFactor2
		rot2.ApplyPush(pushX, pushY)
	}

	// ===== ЛОГИКА ИЗМЕНЕНИЯ УГЛА =====
	resolveAngleCollision(rot1, rot2, hit.Normal, hit.Data, 1.0/60.0)

}

// ============================================================
// ЛОГИКА ИЗМЕНЕНИЯ УГЛА ПРИ СТОЛКНОВЕНИИ ОРУЖИЙ (с ограничением скорости)
// ============================================================

// resolveAngleCollision - изменяет угол оружия и пользователя при столкновении
// Параметры:
//   - weapon1, weapon2: сталкивающиеся оружия
//   - normal: нормаль столкновения
//   - penetration: глубина проникновения
//   - dt: дельта времени

func resolveAngleCollision(weapon1, weapon2 RotatableHitBoxer, normal v.Vec, penetration float64, dt float64) {
	if weapon1 == nil || weapon2 == nil {
		return
	}

	user1 := weapon1.GetWeaponUser()
	user2 := weapon2.GetWeaponUser()
	if user1 == nil || user2 == nil {
		return
	}

	// ===== 1. СРАВНИВАЕМ СИЛЫ ВРАЩЕНИЯ =====
	force1 := user1.GetRotationForce()
	force2 := user2.GetRotationForce()

	// Разница сил (кто доминирует)
	forceDiff := force1 - force2

	// ===== 2. НАПРАВЛЕНИЕ ТОЛЧКА =====
	anglePush1 := math.Atan2(-normal.Y, -normal.X) * 180 / math.Pi
	anglePush2 := math.Atan2(normal.Y, normal.X) * 180 / math.Pi

	// ===== 3. СИЛА ТОЛЧКА =====
	pushForce := 3 * penetration

	// ===== 4. ПРИМЕНЯЕМ ТОЛЧОК К УГЛУ ПОЛЬЗОВАТЕЛЯ =====
	// Доминирующий получает меньше, слабый — больше
	if forceDiff > 0 {
		// Первый доминирует
		user1.ApplyAngularPush(anglePush1 * pushForce * 0.3) // ← слабый толчок
		user2.ApplyAngularPush(anglePush2 * pushForce * 1.5) // ← сильный толчок
	} else if forceDiff < 0 {
		// Второй доминирует
		user1.ApplyAngularPush(anglePush1 * pushForce * 1.5)
		user2.ApplyAngularPush(anglePush2 * pushForce * 0.3)
	} else {
		// Равные силы — оба получают одинаковый толчок
		user1.ApplyAngularPush(anglePush1 * pushForce * 0.8)
		user2.ApplyAngularPush(anglePush2 * pushForce * 0.8)
	}
}

// SAT_OBBvsAABB - SAT для OBB vs AABB
func SAT_OBBvsAABB(
	cx, cy, hw, hh, angle float64,
	px, py, hw2, hh2 float64,
) (bool, v.Vec, float64) {

	rad := angle * math.Pi / 180

	// Оси OBB
	axisOBB_X := v.Vec{X: math.Cos(rad), Y: math.Sin(rad)}
	axisOBB_Y := v.Vec{X: -math.Sin(rad), Y: math.Cos(rad)}

	// Оси AABB
	axisAABB_X := v.Vec{X: 1, Y: 0}
	axisAABB_Y := v.Vec{X: 0, Y: 1}

	axes := []v.Vec{axisOBB_X, axisOBB_Y, axisAABB_X, axisAABB_Y}

	dx := px - cx
	dy := py - cy

	minOverlap := math.MaxFloat64
	var minAxis v.Vec

	for _, axis := range axes {
		// Проекция расстояния
		distProj := math.Abs(dx*axis.X + dy*axis.Y)

		// Проекция OBB
		projOBB := hw*math.Abs(axis.X*math.Cos(rad)+axis.Y*math.Sin(rad)) +
			hh*math.Abs(-axis.X*math.Sin(rad)+axis.Y*math.Cos(rad))

		// Проекция AABB
		projAABB := hw2*math.Abs(axis.X) + hh2*math.Abs(axis.Y)

		overlap := projOBB + projAABB - distProj

		if overlap < 0 {
			return false, v.Vec{}, 0
		}

		if overlap < minOverlap {
			minOverlap = overlap
			minAxis = axis
		}
	}

	// Нормаль — от OBB к AABB
	if dx*minAxis.X+dy*minAxis.Y < 0 {
		minAxis.X = -minAxis.X
		minAxis.Y = -minAxis.Y
	}

	return true, minAxis, minOverlap
}
