package enemy

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// ============================================================
// ROTATION - СТАНДАРТНАЯ ЛОГИКА ВРАЩЕНИЯ ДЛЯ ВРАГА
// ============================================================

// RotationComponent - компонент вращения для врага
// Аналог PlayerHead, но упрощённый
type RotationComponent struct {
	Angle           float64 // текущий угол поворота (градусы)
	AngularVelocity float64 // текущая угловая скорость
	TargetAngle     float64 // целевой угол (куда враг смотрит)

	// Параметры вращения
	RotationSpeed float64 // максимальная скорость вращения
	Acceleration  float64 // ускорение вращения
	Deceleration  float64 // замедление вращения
	Smoothing     float64 // плавность поворота (0.05-0.2)
}

// NewRotationComponent - создаёт компонент вращения с настройками по умолчанию
func NewRotationComponent(rotSpeed, aceleration, deceleration, smoothing float64) *RotationComponent {
	return &RotationComponent{
		Angle:           0.0,
		AngularVelocity: 0.0,
		TargetAngle:     0.0,
		RotationSpeed:   rotSpeed,
		Acceleration:    aceleration,  // ускорение
		Deceleration:    deceleration, // замедление
		Smoothing:       smoothing,    // плавность
	}
}

// UpdateRotation - обновляет вращение врага
// direction: 1 = вправо, -1 = влево, 0 = не вращается
func (r *RotationComponent) UpdateRotation(direction float64, dt float64) {
	// Если направление задано — вращаемся
	if direction != 0 {
		r.AngularVelocity += direction * r.Acceleration * dt
	} else {
		// Замедление
		dec := r.Deceleration * dt
		if math.Abs(r.AngularVelocity) > dec {
			r.AngularVelocity -= math.Copysign(dec, r.AngularVelocity)
		} else {
			r.AngularVelocity = 0
		}
	}

	// Ограничиваем максимальную скорость
	if r.AngularVelocity > r.RotationSpeed {
		r.AngularVelocity = r.RotationSpeed
	}
	if r.AngularVelocity < -r.RotationSpeed {
		r.AngularVelocity = -r.RotationSpeed
	}

	// Применяем вращение
	r.Angle += r.AngularVelocity * dt
	r.TargetAngle = r.Angle

	// Нормализуем угол
	r.Angle = math.Mod(r.Angle, 360)
	if r.Angle < 0 {
		r.Angle += 360
	}
}

// GetAngle - возвращает текущий угол
func (r *RotationComponent) GetAngle() float64 {
	return r.Angle
}

// SetAngle - устанавливает угол
func (r *RotationComponent) SetAngle(angle float64) {
	r.Angle = angle
	r.TargetAngle = angle
}

// SetTargetAngle - устанавливает целевой угол (плавный поворот)
func (r *RotationComponent) SetTargetAngle(target float64) {
	r.TargetAngle = target
}

// UpdateTargetAngle - плавно поворачивается к целевой позиции
func (r *RotationComponent) UpdateTargetAngle(dt float64) {
	diff := r.TargetAngle - r.Angle
	if diff > 180 {
		diff -= 360
	} else if diff < -180 {
		diff += 360
	}

	if math.Abs(diff) < 0.5 {
		r.Angle = r.TargetAngle
		return
	}

	r.Angle += diff * r.Smoothing
	r.Angle = math.Mod(r.Angle, 360)
	if r.Angle < 0 {
		r.Angle += 360
	}
}

// GetDirection - возвращает направление как вектор (для движения)
func (r *RotationComponent) GetDirection() (float64, float64) {
	angleRad := r.Angle * math.Pi / 180
	return math.Cos(angleRad), math.Sin(angleRad)
}

func DrawRotatedRect(screen *ebiten.Image, cx, cy, w, h float64, angleDeg float64, clr color.RGBA) {
	angleReg := angleDeg * math.Pi / 180

	rectImg := ebiten.NewImage(int(w), int(h))
	rectImg.Fill(clr)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(angleReg)
	op.GeoM.Translate(cx, cy)

	screen.DrawImage(rectImg, op)
}
