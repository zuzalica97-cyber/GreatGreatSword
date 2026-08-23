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
	Tag() string
	IsActive() bool
}

type Drawler interface {
	Draw(screen *ebiten.Image, camera *kamera.Camera)
}

// ============================================================
// ВРАГ
// ============================================================

// Enemy - интерфейс для всех врагов
type Enemy interface {
	hitboxes.HitBoxer // для коллизий

	// Базовые параметры
	GetHealth() float64
	SetHealth(health float64)
	GetDamage() float64
	GetPosition() (float64, float64)
	SetPosition(x, y float64)
	GetSpeed() float64
	SetSpeed(speed float64)
	GetSpeedXY() (float64, float64)
	SetSpeedXY(x, y float64)
	GetDirection() (float64, float64)
	SetDirection(x, y float64)
	UpdateDirection(targetX, targetY float64)
	GetTargetDirecton() (targetX, targetY float64)
	SetTargetDirecton(targetX, targetY float64)

	// Действия
	TakeDamage(damage float64) // получить урон
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
