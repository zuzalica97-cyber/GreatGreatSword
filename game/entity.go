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
	Entity                  // Update, Draw, Tag, IsActive
	hitboxes.HitBoxer       // для коллизий
	hitboxes.LetterSender   // для передачи эффектов
	hitboxes.LetterReceiver // для получения эффектов
	hitboxes.EffectUser     // для работы с эффектами

	// Данные врага
	GetHealth() float64
	GetPosition() (float64, float64)
	SetPosition(x, y float64)
	GetSize() int

	// Смерть
	OnDeath(*hitboxes.CollisionManager) // что делать при смерти

}

// ============================================================
// ПУЛЯ / СНАРЯД
// ============================================================

// Projectile - интерфейс для пуль и снарядов
type Projectile interface {
	Entity
	hitboxes.HitBoxer // для коллизий

	GetWeapon() GameWeapon
	GetDamage() int
	GetDirection() (float64, float64)
	GetSpeed() float64
	SetDirection(vx, vy float64)

	// При попадании в цель
	OnHit(target Entity, world WorldView)
}

// ============================================================
// ОРУЖИЕ
// ============================================================

type GameWeapon interface {
	Drawler
	Entity
	hitboxes.HitBoxer
	hitboxes.LetterSender

	IsAttached() bool
	GetPosition() (x, y float64)
	SetPosition(x, y float64)

	// Проризание (насколько глубоко оружие входит в объект)
	GetPenetration() float64
	SetPenetration(penetration float64)
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
	GetWeapon() GameWeapon
}

// ============================================================
// НАНОСИТЕЛЬ УРОНА
// ============================================================

// DamageDealer - объект, который может наносить урон
type DamageDealer interface {
	GetDamage() int
	GetDamageType() string // "melee", "projectile", "explosion"
}
