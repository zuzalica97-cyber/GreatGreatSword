package movmentcommon

import "math"

// SlideMovement - движение со скольжением, зависящее от скорости
// Параметры:
//   - speed: текущая скорость
//   - targetSpeed: желаемая скорость
//   - dirX, dirY: текущее направление
//   - targetDirX, targetDirY: желаемое направление
//   - friction: трение (чем выше, тем дольше скользит) 0.95-0.999
//   - acceleration: ускорение (как быстро разгоняется)
//   - dt: дельта времени
func SlideMovement(
	speed float64,
	targetSpeed float64,
	dirX, dirY float64,
	targetDirX, targetDirY float64,
	friction float64,
	acceleration float64,
	dt float64,
) (float64, float64, float64) {

	// ===== 1. СКОРОСТЬ (ускорение + трение) =====
	// Разница между целевой и текущей скоростью
	diff := targetSpeed - speed

	// Ускорение зависит от разницы: чем больше разрыв, тем быстрее разгон
	if diff > 0 {
		speed += diff * acceleration * dt
	} else if diff < 0 {
		speed -= diff * acceleration * dt // отрицательное ускорение (торможение)
	}

	// Трение (всегда замедляет)
	speed *= (1 - friction*dt)

	// Если скорость очень маленькая — останавливаем
	if math.Abs(speed) < 0.05 {
		speed = 0
	}

	// ===== 2. НАПРАВЛЕНИЕ (зависит от скорости) =====
	// Чем выше скорость, тем сложнее изменить направление
	// Скорость 0 → мгновенный поворот (speedFactor = 1)
	// Скорость 500 → медленный поворот (speedFactor = 0.2)
	speedFactor := 1.0 - (math.Min(math.Abs(speed), 500) / 500 * 1.5)
	if speedFactor < 0.2 {
		speedFactor = 0.2
	}

	// Плавное изменение направления
	newDirX := dirX*(1-speedFactor) + targetDirX*speedFactor
	newDirY := dirY*(1-speedFactor) + targetDirY*speedFactor

	// Нормализуем
	len := math.Sqrt(newDirX*newDirX + newDirY*newDirY)
	if len > 0.01 {
		newDirX /= len
		newDirY /= len
	}

	return speed, newDirX, newDirY
}

// MoveWithInertia - простое движение с инерцией
// Параметры:
//   - speed: текущая скорость
//   - targetSpeed: желаемая скорость (куда хотим прийти)
//   - dirX, dirY: текущее направление
//   - targetDirX, targetDirY: желаемое направление
//   - inertia: инерция (0.9-0.99) — чем выше, тем больше инерции
//   - dt: дельта времени
func MoveWithInertia(
	speed float64,
	targetSpeed float64,
	dirX, dirY float64,
	targetDirX, targetDirY float64,
	inertia float64,
	dt float64,
) (float64, float64, float64) {

	// 1. Скорость плавно стремится к целевой
	speed += (targetSpeed - speed) * inertia * dt * 10

	// 2. Направление плавно стремится к целевому
	newDirX := dirX + (targetDirX-dirX)*inertia*dt*10
	newDirY := dirY + (targetDirY-dirY)*inertia*dt*10

	// Нормализуем направление
	len := math.Sqrt(newDirX*newDirX + newDirY*newDirY)
	if len > 0.01 {
		newDirX /= len
		newDirY /= len
	}

	return speed, newDirX, newDirY
}
