package game

import (
	"great-sword/game/hitboxes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

// ============================================================
// БАЗОВАЯ СУЩНОСТЬ
// ============================================================

type Entity interface {
	Update(world WorldView, manager *hitboxes.CollisionManager) bool
	Draw(screen *ebiten.Image, camera *kamera.Camera)
	Tag() string
	IsActive() bool
}

// ============================================================
// ВРАГ
// ============================================================

// Enemy - интерфейс для всех врагов
type Enemy interface {
	Entity
	hitboxes.HitBoxer // для коллизий

	// Базовые параметры
	GetHealth() int
	SetHealth(health int)
	GetDamage() int
	GetPosition() (float64, float64)
	SetPosition(x, y float64)
	GetSpeed() float64
	SetSpeed(speed float64)
	GetSpeedXY() (float64, float64)
	SetSpeedXY(x, y float64)

	// Действия
	TakeDamage(damage int)   // получить урон
	OnDeath(world WorldView) // что делать при смерти
}

// ============================================================
// ПУЛЯ / СНАРЯД
// ============================================================

// Projectile - интерфейс для пуль и снарядов
type Projectile interface {
	Entity
	hitboxes.HitBoxer // для коллизий

	GetDamage() int
	GetDirection() (float64, float64)
	GetSpeed() float64
	SetDirection(vx, vy float64)

	// При попадании в цель
	OnHit(target Entity, world WorldView)
}

// ============================================================
// ИГРОК
// ============================================================

// Player - интерфейс для игрока
type Player interface {
	Entity
	hitboxes.HitBoxer

	GetHealth() int
	SetHealth(health int)
	TakeDamage(damage int)
	GetPosition() (float64, float64)
	SetPosition(x, y float64)
}

// ============================================================
// НАНОСИТЕЛЬ УРОНА
// ============================================================

// DamageDealer - объект, который может наносить урон
type DamageDealer interface {
	GetDamage() int
	GetDamageType() string // "melee", "projectile", "explosion"
}
