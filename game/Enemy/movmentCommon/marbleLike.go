package movmentcommon

import "math"

// ============================================================
// MarbleMovementParams - параметры движения как в Marble Kingdoms
// ============================================================

type MarbleMovementParams struct {
	// Ускорение (пикселей/сек²)
	Acceleration float64

	// Максимальная скорость (пикселей/сек)
	MaxSpeed float64

	// Трение (коэффициент замедления, 0.95-0.99)
	// Чем ближе к 1, тем дольше катится
	Friction float64

	// Масса (влияет на ускорение и отталкивание)
	Mass float64

	// Дистанция, с которой враг начинает реагировать
	MaxDistance float64

	// Минимальная дистанция (враг пытается держаться)
	MinDistance float64

	// Сила отталкивания от стен (при столкновении)
	BounceForce float64
}

func NewMarbleMovmentParams(
	acceleration, maxSpeed, friction, mass, maxDsitance, minDistance, bounceForce float64,
) *MarbleMovementParams {
	return &MarbleMovementParams{
		Acceleration: acceleration,
		MaxSpeed:     maxSpeed,
		Friction:     friction,
		Mass:         mass,
		MaxDistance:  maxDsitance,
		MinDistance:  minDistance,
		BounceForce:  bounceForce,
	}
}

// ============================================================
// MarbleMovementResult - результат расчёта движения
// ============================================================

type MarbleMovementResult struct {
	// Новые координаты врага
	NewX, NewY float64

	// Новая скорость (для сохранения инерции)
	NewSpeed float64

	// Направление движения (для анимации)
	DirectionX, DirectionY float64

	// Активно ли движение
	IsMoving bool
}

// ============================================================
// ФУНКЦИЯ ПЕРЕДВИЖЕНИЯ КАК В MARBLE KINGDOMS
// ============================================================

// CalculateMarbleMovement - рассчитывает движение шара к цели
// Параметры:
//   - enemyX, enemyY: текущая позиция врага
//   - playerX, playerY: позиция игрока
//   - currentSpeed: текущая скорость врага
//   - dt: дельта времени
//   - params: параметры движения
//
// Возвращает: результат движения (новые координаты, скорость, направление)
func CalculateMarbleMovement(
	enemyX, enemyY float64,
	playerX, playerY float64,
	currentSpeed float64,
	dt float64,
	params MarbleMovementParams,
) MarbleMovementResult {

	result := MarbleMovementResult{
		NewX:     enemyX,
		NewY:     enemyY,
		NewSpeed: currentSpeed,
		IsMoving: false,
	}

	// ===== 1. ВЫЧИСЛЯЕМ НАПРАВЛЕНИЕ К ИГРОКУ =====
	dx := playerX - enemyX
	dy := playerY - enemyY
	distance := math.Sqrt(dx*dx + dy*dy)

	// ===== 2. ЕСЛИ ВРАГ СЛИШКОМ ДАЛЕКО — ОСТАНАВЛИВАЕМСЯ =====
	if distance > params.MaxDistance || distance < params.MinDistance {
		// Применяем трение (замедление)
		result.NewSpeed = currentSpeed * (1 - params.Friction*dt)

		// Если скорость почти нулевая — останавливаем
		if math.Abs(result.NewSpeed) < 0.1 {
			result.NewSpeed = 0
			result.IsMoving = false
			return result
		}

		// Если есть скорость — продолжаем движение по инерции
		if currentSpeed != 0 {
			result.IsMoving = true
			// Сохраняем последнее направление
		}
		return result
	}

	// ===== 3. РАССЧИТЫВАЕМ УСКОРЕНИЕ =====
	// Чем ближе к игроку, тем больше ускорение
	// t = 0 когда distance = MaxDistance (нет ускорения)
	// t = 1 когда distance = MinDistance (максимальное ускорение)
	t := 1.0 - (distance-params.MinDistance)/(params.MaxDistance-params.MinDistance)

	// Ограничиваем t от 0 до 1
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	// Ускорение с учётом массы
	acceleration := params.Acceleration * t / params.Mass

	// ===== 4. НОРМАЛИЗУЕМ НАПРАВЛЕНИЕ =====
	dirX := dx / distance
	dirY := dy / distance

	// ===== 5. ПРИМЕНЯЕМ ФИЗИКУ =====
	// Добавляем ускорение к скорости
	result.NewSpeed = currentSpeed + acceleration*dt

	// Применяем трение
	result.NewSpeed *= (1 - params.Friction*dt)

	// Ограничиваем максимальную скорость
	if result.NewSpeed > params.MaxSpeed {
		result.NewSpeed = params.MaxSpeed
	}
	if result.NewSpeed < -params.MaxSpeed {
		result.NewSpeed = -params.MaxSpeed
	}

	// ===== 6. ВЫЧИСЛЯЕМ НОВУЮ ПОЗИЦИЮ =====
	result.NewX = enemyX + dirX*result.NewSpeed*dt
	result.NewY = enemyY + dirY*result.NewSpeed*dt

	result.DirectionX = dirX
	result.DirectionY = dirY
	result.IsMoving = true

	return result
}
