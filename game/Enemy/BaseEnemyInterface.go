package enemy

import (
	"fmt"
	movmentcommon "great-sword/game/Enemy/movmentCommon"
	"great-sword/game/effects"
	"great-sword/game/hitboxes"

	"image/color"
	"math"
)

var _ hitboxes.EffectUser = (*BaseEnemy)(nil)
var _ hitboxes.HitBoxer = (*BaseEnemy)(nil)
var _ hitboxes.LetterReceiver = (*BaseEnemy)(nil)
var _ hitboxes.LetterSender = (*BaseEnemy)(nil) // проверяем, что реализует интерфейс

// ============================================================
// БАЗОВАЯ СТРУКТУРА ВРАГА
// ============================================================

type BaseEnemy struct {
	// Позиция
	X, Y float64

	// Размер
	Size int

	// Характеристики
	Health   float64
	Damage   int
	Speed    float64
	MaxSpeed float64

	// Состояния
	Active           bool
	CooldownActive   bool
	CooldownTimer    float64
	CooldownDuration float64

	// Цвет
	Color color.RGBA

	// Уникальный ID
	Index int

	// Цель (игрок)
	TargetX, TargetY float64
	TargetDistance   float64

	// Скорость
	SpeedX, SpeedY float64
	CurrentSpeed   float64

	DirectionX, DirectionY float64

	// Вес объекта (чем больше, тем сложнее сдвинуть)
	Weight float64

	// Плотность объекта (0.0 - 1.0)
	// 0.0 = проницаемый, 1.0 = твёрдый (стена)
	Density float64

	// Тег
	TagName string

	HasAuraField        bool // имеет ли ауру
	AffectedByAuraField bool // реагирует ли на чужую ауру

	EffectsManagerEnemy *effects.EffectManager
	Effects             []hitboxes.Effect // список активных эффектов

	// Письма (для системы писем)
	letters []*hitboxes.Letter
}

// ============================================================
// КОНСТРУКТОР
// ============================================================

func NewBaseEnemy(x, y float64, size int, health float64, damage int, speed, maxSpeed float64, color color.RGBA, weight, density float64, tag string) *BaseEnemy {
	return &BaseEnemy{
		X:                   x,
		Y:                   y,
		Size:                size,
		Health:              health,
		Damage:              damage,
		Speed:               speed,
		MaxSpeed:            maxSpeed,
		Active:              true,
		Color:               color,
		Weight:              weight,
		Density:             density,
		TagName:             tag,
		CooldownDuration:    2.0,
		HasAuraField:        true,
		AffectedByAuraField: true,
	}
}

// ============================================================
// УДАЛЁННЫЕ ФУНКЦИИ (теперь в BaseEnemy):
// - GetPosition / SetPosition
// - GetSpeed / SetSpeed
// - GetSpeedXY / SetSpeedXY
// - IsMoving
// - ApplyPush
// - GetTarget / SetTarget / GetTargetDistance
// - GetAABB / GetHitBoxID
// - IsActive / IsStatic
// - Tag
// - GetHealth / SetHealth / TakeDamage
// - MoveTowardsTarget
// - ActivateCooldown / UpdateCooldown
// ============================================================

// ============================================================
// ОСТАВЛЕННЫЕ ФУНКЦИИ (специфичные для этого врага)
// ============================================================

// ============================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С ПОЗИЦИЕЙ
// ============================================================

func (b *BaseEnemy) GetPosition() (float64, float64) {
	return b.X, b.Y
}

func (b *BaseEnemy) SetPosition(x, y float64) {
	b.X = x
	b.Y = y
}

func (b *BaseEnemy) GetSize() int {
	return b.Size
}
func (b *BaseEnemy) SetSize(size int) {
	b.Size = size
}

// ============================================================
// МЕТОДЫ ДЛЯ РАБОТЫ СО СКОРОСТЬЮ
// ============================================================

func (b *BaseEnemy) GetSpeed() float64 {
	return b.CurrentSpeed
}

func (b *BaseEnemy) SetSpeed(speed float64) {
	b.CurrentSpeed = speed
}

func (b *BaseEnemy) GetSpeedXY() (float64, float64) {
	return b.SpeedX, b.SpeedY
}

func (b *BaseEnemy) SetSpeedXY(vx, vy float64) {
	b.SpeedX = vx
	b.SpeedY = vy
}

// ============================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С ЦЕЛЬЮ
// ============================================================

func (b *BaseEnemy) GetTarget() (float64, float64) {
	return b.TargetX, b.TargetY
}

func (b *BaseEnemy) GetTargetDistance() float64 {
	return b.TargetDistance
}

func (b *BaseEnemy) SetTarget(x, y float64) {
	b.TargetX = x
	b.TargetY = y
	dx := x - b.X
	dy := y - b.Y
	b.TargetDistance = math.Sqrt(dx*dx + dy*dy)
}

// ============================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С НАПРАВЛЕНИЕМ
// ============================================================

// GetDirection - возвращает текущее направление движения врага
// Если враг не двигается, возвращает (0, 0)
func (b *BaseEnemy) GetDirection() (float64, float64) {
	return b.DirectionX, b.DirectionY
}

// SetDirection - устанавливает направление движения
func (b *BaseEnemy) SetDirection(x, y float64) {
	b.DirectionX = x
	b.DirectionY = y
}

// UpdateDirection - обновляет направление к цели
// Вычисляет направление от текущей позиции к целевой
func (b *BaseEnemy) UpdateDirection(targetX, targetY float64) {
	dx := targetX - b.X
	dy := targetY - b.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0.01 {
		b.DirectionX = dx / distance
		b.DirectionY = dy / distance
	} else {
		b.DirectionX = 0
		b.DirectionY = 0
	}
}

// GetTargetDirection - возвращает направление к цели (без сохранения)
func (b *BaseEnemy) GetTargetDirection(targetX, targetY float64) (float64, float64) {
	dx := targetX - b.X
	dy := targetY - b.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0.01 {
		return dx / distance, dy / distance
	}
	return 0, 0
}

// ============================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА LetterSender
// ============================================================

// В враге (получатель)
func (b *BaseEnemy) OnCollision(effects []hitboxes.Effect) {
	for _, effect := range effects {
		b.AddEffect(effect)
	}
}

// ============================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА LetterSender (ОТПРАВКА ПИСЕМ)
// ============================================================

// GetEffectsForTransfer - возвращает эффекты для передачи
func (b *BaseEnemy) GetEffectsForTransfer() []hitboxes.Effect {
	var clones []hitboxes.Effect
	for _, letter := range b.letters {
		if !letter.CanDeliver() {
			continue
		}

		letter.Deliver()

		// Возвращаем клоны эффектов для передачи
		for _, effect := range b.Effects {
			clones = append(clones, effect.Clone())
		}
	}
	return clones
}

// CanSendEffects - можно ли отправить эффекты сейчас
func (b *BaseEnemy) CanSendEffects() bool {
	return len(b.letters) > 0
}

// OnEffectsSent - вызывается после отправки эффектов
func (b *BaseEnemy) OnEffectsSent() {
	// После отправки очищаем эффекты (или можно оставить)
}

// ============================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С ЭФФЕКТАМИ
// ============================================================

func (b *BaseEnemy) AddEffect(effect hitboxes.Effect) { //Нужно реализовать OnColision прямо в BaseEnemy чтобы всё работало.

	// Проверяем, есть ли уже такой эффект
	for _, e := range b.Effects {
		if e.GetType() == effect.GetType() && e.IsActive() && effect.CanExtend() {
			// Суммируем время
			e.SetDuration(e.GetDuration() + effect.GetDuration()*0.5)
			e.SetTimer(e.GetDuration() * 0.5)
			fmt.Printf("Суммируем эффект: %s, новое время: %.1f\n",
				effect.GetType(), e.GetDuration())
			return
		}
	}

	b.Effects = append(b.Effects, effect)
	effect.Apply(b)
}

func (b *BaseEnemy) RemoveEffect(effectID string) {
	for i, e := range b.Effects {
		if e.GetID() == effectID {
			b.Effects = append(b.Effects[:i], b.Effects[i+1:]...)
			return
		}
	}
}

func (b *BaseEnemy) HasEffect(effectType string) bool {
	for _, e := range b.Effects {
		if e.GetType() == effectType && e.IsActive() {
			return true
		}
	}
	return false
}

// UpdateEffects - обновляет все эффекты врага
// Возвращает true, если есть активные эффекты
func (b *BaseEnemy) UpdateEffects(dt float64) bool {
	hasActive := false
	for i := 0; i < len(b.Effects); i++ {
		effect := b.Effects[i]
		if effect.Update(b, dt) || !effect.IsActive() {
			// Эффект завершён — удаляем
			b.Effects = append(b.Effects[:i], b.Effects[i+1:]...)
			i--
		} else {
			hasActive = true
		}
	}
	return hasActive
}

// ============================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА EffectUser (hitboxes.EffectUser)
// ============================================================

func (b *BaseEnemy) GetEffects() []hitboxes.Effect {
	return b.Effects
}

func (b *BaseEnemy) ClearEffects() {
	b.Effects = make([]hitboxes.Effect, 0)
}

func (b *BaseEnemy) GetMaxSpeed() float64 {
	return b.MaxSpeed
}

// ============================================================
// МЕТОДЫ ДЛЯ РАБОТЫ СО ЗДОРОВЬЕМ
// ============================================================

func (b *BaseEnemy) GetHealth() float64 {
	return b.Health
}

func (b *BaseEnemy) SetHealth(health float64) {
	b.Health = health
	if b.Health <= 0 {
		b.Active = false
	}
}

func (b *BaseEnemy) TakeDamage(damage float64) {
	b.Health -= damage
	if b.Health <= 0 {
		b.Active = false
	}
}

func (b *BaseEnemy) IsActive() bool {
	return b.Active
}

// ============================================================
//МЕТОДЫ ДЛЯ КУЛДАУНА
// ============================================================

func (b *BaseEnemy) ActivateCooldown() {
	b.CooldownActive = true
	b.CooldownTimer = b.CooldownDuration
}

func (b *BaseEnemy) UpdateCooldown(dt float64) {
	if b.CooldownActive {
		b.CooldownTimer -= dt
		if b.CooldownTimer <= 0 {
			b.CooldownActive = false
		}
	}
}

// ============================================================
// МЕТОДЫ ДЛЯ ОТТАЛКИВАНИЯ
// ============================================================

func (b *BaseEnemy) ApplyPush(x, y float64) {
	b.X += x
	b.Y += y
}

// ============================================================
// РЕАЛИЗАЦИЯ hitboxes.HitBoxer
// ============================================================

func (b *BaseEnemy) GetAABB() (posX, posY, halfW, halfH float64) {
	halfSize := float64(b.Size) / 2
	return b.X + halfSize, b.Y + halfSize, halfSize, halfSize
}

func (b *BaseEnemy) GetHitBoxID() string {
	return b.TagName + "_" + string(rune(b.Index))
}

func (b *BaseEnemy) IsStatic() bool {
	return false
}

// GetWeight - возвращает вес объекта
func (b *BaseEnemy) GetWeight() float64 {
	return b.Weight
}

// GetDensity - возвращает плотность объекта
func (b *BaseEnemy) GetDensity() float64 {
	return b.Density
}

func (b *BaseEnemy) HasAura() bool {
	return b.HasAuraField
}

func (b *BaseEnemy) AffectedByAura() bool {
	return b.AffectedByAuraField
}

// ============================================================
// ИНТЕРФЕЙС game.Entity
// ============================================================

func (b *BaseEnemy) Tag() string {
	return b.TagName
}

// ============================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ============================================================

// IsMoving - проверяет, движется ли враг в данный момент
// Возвращает true, если есть скорость по X, Y или общая скорость
func (b *BaseEnemy) IsMoving() bool {
	return b.SpeedX != 0 || b.SpeedY != 0 || b.CurrentSpeed != 0
}

// GetDistToTarget - вычисляет расстояние и вектор до цели (игрока)
// Возвращает:
//   - dist: расстояние до цели
//   - dx, dy: вектор от врага к цели (не нормализованный)
func (b *BaseEnemy) GetDistToTarget() (dist float64, dx float64, dy float64) {
	X, Y := b.GetTarget()

	dx = X - b.X
	dy = Y - b.Y
	dist = math.Sqrt(dx*dx + dy*dy)
	return dist, dx, dy
}

// GetTargetDir - возвращает нормализованное направление к цели (игроку)
// Возвращает:
//   - targetDirX, targetDirY: единичный вектор направления к игроку
//   - если цель не достижима, возвращает (0, 0)
func (b *BaseEnemy) GetTargetDir() (float64, float64) {
	var targetDirX, targetDirY float64

	dist, dx, dy := b.GetDistToTarget()

	if dist > 0.01 {
		targetDirX = dx / dist
		targetDirY = dy / dist
	} else {
		targetDirX, targetDirY = 0, 0
	}
	return targetDirX, targetDirY
}

// EnemySlideMovmentFunc - применяет движение со скольжением к врагу
// Параметры:
//   - friction: коэффициент трения (0.95-0.999) — чем выше, тем дольше скользит
//   - acceleration: скорость разгона (чем выше, тем быстрее разгоняется)
//   - dt: дельта времени (обычно 1/60)
//
// Возвращает:
//   - newSpeed: новая скорость врага
//   - newDirX, newDirY: новое направление движения
func (b *BaseEnemy) EnemySlideMovmentFunc(friction, acceleration, dt float64) (float64, float64, float64) {
	currentDirX, currentDirY := b.GetDirection()

	targetDirX, targetDirY := b.GetTargetDir()

	newSpeed, newDirX, newDirY := movmentcommon.SlideMovement(
		b.CurrentSpeed,
		b.Speed,
		currentDirX, currentDirY,
		targetDirX, targetDirY,
		friction,
		acceleration,
		dt,
	)
	return newSpeed, newDirX, newDirY
}
