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

	weight1 := weapon1.GetWeight()
	weight2 := weapon2.GetWeight()

	// ===== 1. НАПРАВЛЕНИЕ ТОЛЧКА =====
	anglePush1 := math.Atan2(-normal.Y, -normal.X) * 180 / math.Pi //ДЗ сделай силу вращения и чтобы мечь мог отталкивать врагов а не просто проходить с квозь.
	anglePush2 := math.Atan2(normal.Y, normal.X) * 180 / math.Pi

	// ===== 2. СИЛА ТОЛЧКА =====
	pushForce := 5.0 * penetration

	// ===== 3. МАКСИМАЛЬНАЯ СКОРОСТЬ ИЗМЕНЕНИЯ УГЛА =====
	// Максимальное изменение угла за кадр (градусов/сек)
	maxAngularSpeed := 720.0 // 180 градусов в секунду
	maxAngleChange := maxAngularSpeed * dt

	// ===== 4. ТОЛКАЕМ ПОЛЬЗОВАТЕЛЯ ПЕРВОГО ОРУЖИЯ =====
	user1 := weapon1.GetWeaponUser()
	if user1 != nil {
		force1 := pushForce * (weight1 / (weight1 + 1))
		currentUserAngle := user1.GetAngle()

		// Желаемое изменение угла
		desiredChange := anglePush1 * force1

		// ===== ОГРАНИЧЕНИЕ СКОРОСТИ =====
		if desiredChange > maxAngleChange {
			desiredChange = maxAngleChange
		} else if desiredChange < -maxAngleChange {
			desiredChange = -maxAngleChange
		}

		newUserAngle := currentUserAngle + desiredChange
		newUserAngle = math.Mod(newUserAngle, 360)
		if newUserAngle < 0 {
			newUserAngle += 360
		}

		user1.SetAngle(newUserAngle)
	}

	// ===== 5. ТОЛКАЕМ ПОЛЬЗОВАТЕЛЯ ВТОРОГО ОРУЖИЯ =====
	user2 := weapon2.GetWeaponUser()
	if user2 != nil {
		force2 := pushForce * (weight2 / (weight2 + 1))
		currentUserAngle := user2.GetAngle()

		desiredChange := anglePush2 * force2

		// ===== ОГРАНИЧЕНИЕ СКОРОСТИ =====
		if desiredChange > maxAngleChange {
			desiredChange = maxAngleChange
		} else if desiredChange < -maxAngleChange {
			desiredChange = -maxAngleChange
		}

		newUserAngle := currentUserAngle + desiredChange
		newUserAngle = math.Mod(newUserAngle, 360)
		if newUserAngle < 0 {
			newUserAngle += 360
		}

		user2.SetAngle(newUserAngle)
	}
}
