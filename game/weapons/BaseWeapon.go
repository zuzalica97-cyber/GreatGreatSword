package weapons

import (
	"fmt"
	"great-sword/game"
	"great-sword/game/hitboxes"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

var _ hitboxes.HitBoxer = (*BaseWeapon)(nil) //ДЗ разобратся с частами меча их работой а потом переходить к фабрике монстров
var _ hitboxes.RotatableHitBoxer = (*WeaponPart)(nil)
var _ hitboxes.LetterSender = (*WeaponPart)(nil)
var _ hitboxes.PartRot = (*WeaponPart)(nil)

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
	// Сторона привязки к пользователю
	// "right", "left", "up", "down"
	AttachSide string

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

	// Части оружия
	Parts []*WeaponPart

	// Физические параметры (встроены в структуру)
	Physics WeaponPhysics

	// Письма для эффектов
	Letters []*hitboxes.Letter

	DebugContactX, DebugContactY float64 // точка контакта для отладки
	DebugContactActive           bool    // есть ли столкновение
}

// WeaponPart - часть оружия (отдельная фигура)
type WeaponPart struct {
	// Смещение относительно центра оружия
	OffsetX, OffsetY float64
	// Размер
	Width, Height float64
	// Угол относительно оружия
	AngleOffset float64
	// ID для коллизий
	ID string
	// Письма для эффектов
	Letters []*hitboxes.Letter
	// Вес и плотность
	Weight, Density float64
	// мировые координаты (вычисляются автоматически)
	WorldX, WorldY float64
	WorldAngle     float64

	Collor color.Color

	ParentWeapon *BaseWeapon // ← ссылка на родительское оружие
}

// ============================================================
// КОНСТРУКТОР
// ============================================================

// NewBaseWeapon - создаёт базовое оружие
func NewBaseWeapon(width, height, density, penetration float64, physics WeaponPhysics, side string, id string, user WeaponUser) *BaseWeapon {
	return &BaseWeapon{
		Width:       width,
		Height:      height,
		Density:     density,
		Penetration: penetration,
		Physics:     physics,
		AttachSide:  side,
		Active:      true,
		attached:    false,
		ID:          id,
		MaxLength:   50,
		user:        user,
	}
}

func NewWeaponPart(offsetX, offsetY, width, height, angleOffset float64, id string, collor color.Color, rutWeapon *BaseWeapon) *WeaponPart {
	return &WeaponPart{
		OffsetX:      offsetX,
		OffsetY:      offsetY,
		Width:        width,
		Height:       height,
		AngleOffset:  angleOffset,
		ID:           id,
		Weight:       1.0,
		Density:      0.8,
		Collor:       collor,
		ParentWeapon: rutWeapon,
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
// РЕАЛИЗАЦИЯ ЛОГИКИ ДЛЯ ЧАСТЕЙ ОРУЖИЯ
// ============================================================

// UpdateParts - обновляет мировые координаты всех частей оружия
func (b *BaseWeapon) UpdateParts(dt float64) {
	angleRad := b.Angle * math.Pi / 180

	for _, part := range b.Parts {

		for _, letter := range part.Letters {
			letter.UpdateCoolDown(dt)

		}

		// Поворачиваем смещение части на угол оружия
		rotatedX := part.OffsetX*math.Cos(angleRad) - part.OffsetY*math.Sin(angleRad)
		rotatedY := part.OffsetX*math.Sin(angleRad) + part.OffsetY*math.Cos(angleRad)

		// Мировые координаты части
		part.WorldX = b.PositionX + rotatedX
		part.WorldY = b.PositionY + rotatedY

		// Мировой угол части
		part.WorldAngle = b.Angle + part.AngleOffset
	}
}

func (p *WeaponPart) GetAABB() (posX, posY, halfW, halfH float64) {
	return p.WorldX, p.WorldY, p.Width / 2, p.Height / 2
}

func (p *WeaponPart) GetOBB() (centerX, centerY, halfW, halfH, angle float64) {
	return p.WorldX, p.WorldY, p.Width / 2, p.Height / 2, p.WorldAngle * math.Pi / 180
}

func (p *WeaponPart) GetHitBoxID() string {
	return p.ID
}

func (p *WeaponPart) IsActive() bool {
	return true
}

func (p *WeaponPart) IsStatic() bool {
	return false
}

func (w *WeaponPart) ApplyPush(x, y float64) {
	if w.ParentWeapon != nil {
		w.ParentWeapon.ApplyPush(x, y)
	}
}

func (p *WeaponPart) OnCollision(other hitboxes.HitBoxer) {
}

func (p *WeaponPart) GetWeight() float64 {
	return p.Weight
}

func (p *WeaponPart) GetDensity() float64 {
	return p.Density
}

func (p *WeaponPart) HasAura() bool {
	return false
}

func (p *WeaponPart) AffectedByAura() bool {
	return false
}

// GetAngle - возвращает текущий мировой угол части
func (p *WeaponPart) GetAngle() float64 {
	return p.WorldAngle
}

// GetAngle - возвращает текущий мировой угол части
func (p *WeaponPart) SetAngle(angle float64) {
	p.WorldAngle = angle
}

func (p *WeaponPart) GetTargetAngle() float64 {
	return p.WorldAngle
}

func (p *WeaponPart) GetWeaponUser() hitboxes.WeaponUser {
	if p.ParentWeapon != nil {
		return p.ParentWeapon.GetWeaponUser()
	}
	return nil
}

func (p *WeaponPart) GetRotWeapon() hitboxes.RotatableHitBoxer {
	if p.ParentWeapon != nil {
		return p.ParentWeapon
	}
	return nil
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

	// ===== ВЫБОР СТОРОНЫ =====
	var offsetX, offsetY float64
	var angleOffset float64 // ← смещение угла

	switch b.AttachSide {
	case "right":
		offsetX = math.Cos(angleRad) * distanceFromCenter
		offsetY = math.Sin(angleRad) * distanceFromCenter
		angleOffset = 0
	case "left":
		offsetX = -math.Cos(angleRad) * distanceFromCenter
		offsetY = -math.Sin(angleRad) * distanceFromCenter
		angleOffset = 180
	case "up":
		offsetX = -math.Sin(angleRad) * distanceFromCenter
		offsetY = math.Cos(angleRad) * distanceFromCenter
		angleOffset = -90
	case "down":
		offsetX = math.Sin(angleRad) * distanceFromCenter
		offsetY = -math.Cos(angleRad) * distanceFromCenter
		angleOffset = 90
	default:
		offsetX = math.Cos(angleRad) * distanceFromCenter
		offsetY = math.Sin(angleRad) * distanceFromCenter
		angleOffset = 0
	}

	b.TargetX = centrX + offsetX
	b.TargetY = centrY + offsetY

	// ===== ЦЕЛЕВОЙ УГОЛ С УЧЁТОМ СТОРОНЫ =====
	b.TargetAngle = userAngle + angleOffset
	b.TargetAngle = math.Mod(b.TargetAngle, 360)
	if b.TargetAngle < 0 {
		b.TargetAngle += 360
	}
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

// RegisterParts - добавляет все части оружия в менеджер коллизий
func (b *BaseWeapon) RegisterParts(manager *hitboxes.CollisionManager) {
	for _, part := range b.Parts {
		manager.AddObject(part)
	}
}

// RemoveParts - удаляет все части оружия из менеджера коллизий
func (b *BaseWeapon) RemoveParts(manager *hitboxes.CollisionManager) {
	for _, part := range b.Parts {
		manager.RemoveObject(part)
	}
}

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

// GetEffectsForTransfer - возвращает эффекты для передачи
func (w *WeaponPart) GetEffectsForTransfer(target any) []hitboxes.Effect {
	var effects []hitboxes.Effect
	for _, letter := range w.Letters {
		if !letter.WhiteListLetters(target) || !letter.CanDeliver(target) {
			fmt.Println(7)
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
func (w *WeaponPart) CanSendEffects(target any) bool {
	for _, letter := range w.Letters {
		if letter.CanDeliver(target) && letter.WhiteListLetters(target) {
			return true
		}

	}
	return false
}

//============================================================================
//ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
//============================================================================

func (b *BaseWeapon) StandartDraw(screen *ebiten.Image, camera *kamera.Camera) {
	Color := color.RGBA{0, 100, 200, 255}

	// Создаём временное изображение для меча
	swordImg := ebiten.NewImage(int(b.Width), int(b.Height))
	swordImg.Fill(Color)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-b.Width/2, -b.Height/2)
	op.GeoM.Rotate(b.Angle * math.Pi / 180)
	op.GeoM.Translate(b.PositionX, b.PositionY) // мировые координаты

	// Камера сама применит смещение
	camera.Draw(swordImg, op, screen)

	// ===== ОТРИСОВКА ЧАСТЕЙ =====
	for _, part := range b.Parts {
		partImg := ebiten.NewImage(int(part.Width), int(part.Height))
		partImg.Fill(part.Collor)

		partOp := &ebiten.DrawImageOptions{}
		partOp.GeoM.Translate(-part.Width/2, -part.Height/2)
		partOp.GeoM.Rotate(part.WorldAngle * math.Pi / 180)
		partOp.GeoM.Translate(part.WorldX, part.WorldY)

		camera.Draw(partImg, partOp, screen)
	}

	// Отрисовка точки контакта
	if b.DebugContactActive {
		size := 26.0
		screenX := float32(b.DebugContactX - camera.X - size/2)
		screenY := float32(b.DebugContactY - camera.Y - size/2)

		vector.FillRect(
			screen,
			screenX, screenY,
			float32(size), float32(size),
			color.RGBA{255, 0, 0, 255}, // красный
			true,
		)
		b.DebugContactActive = false // сбрасываем после отрисовки
	}
}

func (b *BaseWeapon) StandartUpdate() {
	dt := 1.0 / 60.0

	for _, letter := range b.Letters {
		letter.UpdateCoolDown(dt)
	}

	// 1. Обновляем целевую позицию (для ориентации)
	b.UpdateAttachmentTarget()

	// 2. Применяем физику "гвоздя" с инерцией
	b.UpdateWeaponAngle(dt)

	b.UpdateWeaponPosition(dt)

	// 4. Обновляем мировые координаты частей
	b.UpdateParts(dt)

}
