package weapons

import (
	"great-sword/game"
	"great-sword/game/hitboxes"
	"math"
)

var _ hitboxes.HitBoxer = (*BaseWeapon)(nil) //ДЗ разобратся с функциями меча и понять его физику. после сделать с помощю baseWeapon  какойнибуть мечь и прикрутитье его к игроку. потом сделать не достающию функция в hitbox

// ============================================================
// BaseWeapon - базовая структура оружия (как у меча)
// ============================================================

type BaseWeapon struct {
	// Привязка к сущности
	user     WeaponUser
	attached bool

	// Позиция и угол
	PositionX, PositionY float64
	TargetX, TargetY     float64
	Angle                float64
	TargetAngle          float64
	AngularVelocity      float64

	AnchorX, AnchorY     float64 // точка крепления (центр пользователя)
	MaxLength            float64 // максимальная длина "нити"
	VelocityX, VelocityY float64 // скорость оружия (для инерции)

	// Размеры
	Width, Height float64

	// Проризание (0.0 - 1.0) - насколько глубоко оружие входит в объект
	Penetration float64

	// Коллизия и состояние
	Active  bool
	Density float64

	// Уникальный ID
	ID string

	// Физические параметры (встроены в структуру)
	Physics WeaponPhysics

	// Письма для эффектов
	Letters []*hitboxes.Letter

	DebugContactX, DebugContactY float64 // точка контакта для отладки
	DebugContactActive           bool    // есть ли столкновение
}

// ============================================================
// КОНСТРУКТОР
// ============================================================

// NewBaseWeapon - создаёт базовое оружие
func NewBaseWeapon(width, height, density, penetration float64, physics WeaponPhysics, id string, user WeaponUser) *BaseWeapon {
	return &BaseWeapon{
		Width:       width,
		Height:      height,
		Density:     density,
		Penetration: penetration,
		Physics:     physics,
		Active:      true,
		attached:    false,
		ID:          id,
		MaxLength:   50,
		user:        user,
	}
}

func (w *BaseWeapon) Debag(x, y float64) {
	w.DebugContactX = x
	w.DebugContactY = y
	w.DebugContactActive = true
}

// ============================================================
// РЕАЛИЗАЦИЯ game.Weapon
// ============================================================

// Attach - привязывает оружие к сущности
func (w *BaseWeapon) Attach(user WeaponUser) {
	w.user = user
	w.attached = true
}

// GetUser - возвращает сущность, к которой привязано оружие
func (w *BaseWeapon) GetUser() WeaponUser {
	return w.user
}

func (w *BaseWeapon) GetWeaponUser() hitboxes.WeaponUser {
	return w.user
}

// IsAttached - проверяет, привязано ли оружие
func (w *BaseWeapon) IsAttached() bool {
	return w.attached
}

// GetPenetration - возвращает уровень проризания оружия
func (w *BaseWeapon) GetPenetration() float64 {
	return w.Penetration
}

// SetPenetration - устанавливает уровень проризания оружия
func (w *BaseWeapon) SetPenetration(penetration float64) {
	w.Penetration = penetration
}

// GetPosition - возвращает позицию оружия
func (w *BaseWeapon) GetPosition() (float64, float64) {
	return w.PositionX, w.PositionY
}

// SetPosition - устанавливает позицию оружия
func (w *BaseWeapon) SetPosition(x, y float64) {
	w.PositionX = x
	w.PositionY = y
}

// GetAngle - возвращает угол поворота оружия
func (w *BaseWeapon) GetAngle() float64 {
	return w.Angle
}

// SetAngle - устанавливает угол поворота оружия
func (w *BaseWeapon) SetAngle(angle float64) {
	w.Angle = angle
}

// ============================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С ФИЗИКОЙ
// ============================================================

// GetPhysics - возвращает физические параметры оружия
func (w *BaseWeapon) GetPhysics() WeaponPhysics {
	return w.Physics
}

// SetPhysics - устанавливает физические параметры оружия
func (w *BaseWeapon) SetPhysics(physics WeaponPhysics) {
	w.Physics = physics
}

// ============================================================
// СТАНДАРТНАЯ ФИЗИКА ОРУЖИЯ (с использованием Physics)
// ============================================================

// UpdateAttachmentTarget - вычисляет целевую позицию оружия относительно пользователя
func (b *BaseWeapon) UpdateAttachmentTarget() {
	if b.user == nil {
		return
	}

	userPosX, userPosY := b.user.GetPosition()

	userSize := b.user.GetSize()
	userAngle := b.user.GetAngle()

	centrX := userPosX + float64(userSize)/2
	centrY := userPosY + float64(userSize)/2

	_, ok := b.user.(game.PlayerLegInter) // костыль

	if ok {
		centrX = userPosX + float64(userSize)
		centrY = userPosY + float64(userSize)

	}

	angleRad := userAngle * math.Pi / 180

	// Используем Physics.Weight для расстояния от центра
	distanceFromCenter := (float64(userSize) / 2) * 1.3

	offsetX := math.Cos(angleRad) * distanceFromCenter
	offsetY := math.Sin(angleRad) * distanceFromCenter

	b.TargetX = centrX + offsetX
	b.TargetY = centrY + offsetY

	b.TargetAngle = userAngle
}

// ============================================================
// ЖЁСТКАЯ ПРИВЯЗКА ОРУЖИЯ (как "гвоздь" в Marble Kingdoms)
// ============================================================
// UpdateWeaponAngle - обновляет угол оружия с плавным следованием и ограничением
func (b *BaseWeapon) UpdateWeaponAngle(dt float64) {
	if b.user == nil {
		return
	}

	maxDeviation := 90.0

	angleDiff := b.TargetAngle - b.Angle
	if angleDiff > 180 {
		angleDiff -= 360
	} else if angleDiff < -180 {
		angleDiff += 360
	}

	//усиления стремелния гдрауса к циливой

	threshold := 5.0

	// Минимальный и максимальный smoothing
	minSmoothing := 0.05
	maxSmoothing := 0.5

	// Вычисляем текущий smoothing в зависимости от отклонения
	absDiff := math.Abs(angleDiff)
	if absDiff <= threshold {
		// Малое отклонение — слабый возврат
		b.Physics.SmoothingReturn = minSmoothing
	} else {
		// Большое отклонение — сильный возврат
		// Чем больше отклонение, тем ближе к maxSmoothing
		t := (absDiff - threshold) / (maxDeviation - threshold)
		if t > 1.0 {
			t = 1.0
		}
		b.Physics.SmoothingReturn = minSmoothing + t*(maxSmoothing-minSmoothing)
	}

	b.Angle += angleDiff * b.Physics.SmoothingReturn

	angleDiff = b.TargetAngle - b.Angle
	if angleDiff > 180 {
		angleDiff -= 360
	} else if angleDiff < -180 {
		angleDiff += 360
	}

	if angleDiff > maxDeviation {
		b.Angle = b.TargetAngle - maxDeviation
	} else if angleDiff < -maxDeviation {
		b.Angle = b.TargetAngle + maxDeviation
	}

	b.Angle = math.Mod(b.Angle, 360)
	if b.Angle < 0 {
		b.Angle += 360
	}
}

// UpdateWeaponPosition - обновляет позицию оружия (вращение вокруг точки соединения)
func (b *BaseWeapon) UpdateWeaponPosition(dt float64) {
	if b.user == nil {
		return
	}

	// Точка соединения (петля)
	anchorX := b.TargetX
	anchorY := b.TargetY

	// Вычисляем позицию центра меча от точки соединения
	angleRad := b.Angle * math.Pi / 180
	halfWidth := b.Width / 2

	offsetX := math.Cos(angleRad) * halfWidth
	offsetY := math.Sin(angleRad) * halfWidth

	targetCenterX := anchorX + offsetX
	targetCenterY := anchorY + offsetY

	// Плавное приближение к целевой позиции
	maxOffset := 7.0
	b.PositionX += (targetCenterX - b.PositionX) * b.Physics.SmoothingPos
	b.PositionY += (targetCenterY - b.PositionY) * b.Physics.SmoothingPos

	// Ограничение позиции
	dx := targetCenterX - b.PositionX
	dy := targetCenterY - b.PositionY
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > maxOffset {
		if dist > 0.01 {
			b.PositionX = targetCenterX - (dx/dist)*maxOffset
			b.PositionY = targetCenterY - (dy/dist)*maxOffset
		}
	}
}

// ============================================================
// РЕАЛИЗАЦИЯ hitboxes.HitBoxer
// ============================================================

// GetAABB - возвращает Axis-Aligned Bounding Box для коллизий
func (w *BaseWeapon) GetAABB() (posX, posY, halfW, halfH float64) {
	return w.PositionX, w.PositionY, w.Width / 2, w.Height / 2
}

// GetHitBoxID - возвращает уникальный ID для системы коллизий
func (w *BaseWeapon) GetHitBoxID() string {
	return w.ID
}

// IsActive - проверяет, активно ли оружие
func (w *BaseWeapon) IsActive() bool {
	return w.Active
}

// IsStatic - всегда false (оружие двигается)
func (w *BaseWeapon) IsStatic() bool {
	return false
}

// ApplyPush - применяет силу отталкивания к оружию и пользователю
func (w *BaseWeapon) ApplyPush(x, y float64) {
	w.PositionX += x
	w.PositionY += y
	if w.user != nil {
		w.user.ApplyPush(x, y)
	}
}

func (w *BaseWeapon) GetTargetAngle() float64 {
	return w.TargetAngle
}

// OnCollision - вызывается при столкновении с другим объектом
func (w *BaseWeapon) OnCollision(other hitboxes.HitBoxer) {
	if w.user != nil {
		tx, ty, _, _ := other.GetAABB()
		wx, wy, _, _ := w.GetAABB()
		dx, dy := tx-wx, ty-wy
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0.01 {
			// Используем ImpactForce и Weight для расчёта силы
			force := w.Physics.ImpactForce / (other.GetWeight() + w.Physics.Weight)
			w.user.ApplyPush((dx/dist)*force, (dy/dist)*force)
		}
	}
}

// GetWeight - возвращает вес оружия
func (w *BaseWeapon) GetWeight() float64 {
	return w.Physics.Weight
}

// GetDensity - возвращает плотность оружия
func (w *BaseWeapon) GetDensity() float64 {
	return w.Density
}

// HasAura - имеет ли оружие ауру отталкивания
func (w *BaseWeapon) HasAura() bool {
	return false
}

// AffectedByAura - реагирует ли оружие на чужую ауру
func (w *BaseWeapon) AffectedByAura() bool {
	return false
}

// Для OBB (точный)
func (b *BaseWeapon) GetOBB() (centerX, centerY, halfW, halfH, angle float64) {
	return b.PositionX, b.PositionY,
		b.Width / 2,
		b.Height / 2,
		b.Angle * math.Pi / 180
}

// ============================================================
// РЕАЛИЗАЦИЯ hitboxes.LetterSender
// ============================================================

// GetEffectsForTransfer - возвращает эффекты для передачи
func (w *BaseWeapon) GetEffectsForTransfer(target any) []hitboxes.Effect {
	var effects []hitboxes.Effect
	for _, letter := range w.Letters {
		if !letter.WhiteListLetters(target) || !letter.CanDeliver(target) {
			continue
		}
		letter.Deliver(target)
		for _, effect := range letter.Effects {
			effects = append(effects, effect.Clone())
		}
	}
	return effects
}

// CanSendEffects - можно ли отправить эффекты
func (w *BaseWeapon) CanSendEffects(target any) bool {
	for _, letter := range w.Letters {
		if letter.CanDeliver(target) && letter.WhiteListLetters(target) {
			return true
		}
	}
	return false
}
