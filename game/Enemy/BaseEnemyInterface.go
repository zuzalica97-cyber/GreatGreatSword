package enemy

import (
	"image/color"
	"math"
)

// ============================================================
// БАЗОВАЯ СТРУКТУРА ВРАГА
// ============================================================

type BaseEnemy struct {
	// Позиция
	X, Y float64

	// Размер
	Size int

	// Характеристики
	Health   int
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

	// Тег
	TagName string
}

// ============================================================
// КОНСТРУКТОР
// ============================================================

func NewBaseEnemy(x, y float64, size, health, damage int, speed, maxSpeed float64, color color.RGBA, tag string) *BaseEnemy {
	return &BaseEnemy{
		X:                x,
		Y:                y,
		Size:             size,
		Health:           health,
		Damage:           damage,
		Speed:            speed,
		MaxSpeed:         maxSpeed,
		Active:           true,
		Color:            color,
		TagName:          tag,
		CooldownDuration: 2.0,
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
// МЕТОДЫ ДЛЯ РАБОТЫ СО ЗДОРОВЬЕМ
// ============================================================

func (b *BaseEnemy) GetHealth() int {
	return b.Health
}

func (b *BaseEnemy) SetHealth(health int) {
	b.Health = health
	if b.Health <= 0 {
		b.Active = false
	}
}

func (b *BaseEnemy) TakeDamage(damage int) {
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

// ============================================================
// ИНТЕРФЕЙС game.Entity
// ============================================================

func (b *BaseEnemy) Tag() string {
	return b.TagName
}

// ============================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ============================================================

func (b *BaseEnemy) IsMoving() bool {
	return b.SpeedX != 0 || b.SpeedY != 0 || b.CurrentSpeed != 0
}
