package movmentcommon

import "math"

// ============================================================
// ClampSpeedDifference - ограничивает разрыв между скоростями
// ============================================================
// Параметры:
//   - previousSpeed: скорость на прошлом кадре
//   - currentSpeed: текущая скорость (желаемая)
//   - maxDiff: максимальный разрыв между скоростями
//
// Возвращает:
//   - newSpeed: скорректированная скорость
//   - isClamped: был ли применён лимит (true если скорость была изменена)
//
// Пример:
//   previousSpeed = 5, currentSpeed = 20, maxDiff = 5
//   → newSpeed = 10 (5 + 5)
//
//   previousSpeed = 5, currentSpeed = 8, maxDiff = 5
//   → newSpeed = 8 (разрыв меньше 5, скорость не меняется)
// ============================================================

func ClampSpeedDiff(prev, curr, maxDiff float64) float64 {
	diff := curr - prev
	if math.Abs(diff) <= maxDiff {
		return curr
	}
	if diff > 0 {
		return prev + maxDiff
	}
	return prev - maxDiff
}
