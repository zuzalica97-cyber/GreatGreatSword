package weapons

import (
	"great-sword/game"
	"great-sword/game/hitboxes"
)

// WeaponUser - сущность, которая может использовать оружие
type WeaponUser interface {
	hitboxes.WeaponUser
	GetWeapon() Weapon
	SetWeapon(weapon Weapon)
	ApplyPush(x, y float64) // для передачи сдвига от оружия
	GetPosition() (float64, float64)
	GetSize() int
	GetSpeedXY() (float64, float64)
	SetSpeedXY(float64, float64)
	GetRotationForce() float64        // текущая сила вращения
	SetRotationForce(force float64)   // установить силу
	GetRotationDirection() float64    // -1 = влево, 0 = стоит, 1 = вправо
	SetRotationDirection(dir float64) // задать направление
	ApplyAngularPush(push float64)    // толчок угла (при столкновении)
}

// Weapon - интерфейс оружия
type Weapon interface {
	game.Drawler
	game.Entity
	hitboxes.HitBoxer
	hitboxes.LetterSender

	// Привязка к пользователю
	Attach(user WeaponUser)
	GetUser() WeaponUser
	IsAttached() bool
	GetPosition() (x, y float64)
	SetPosition(x, y float64)

	// Физика оружия (как в Marble Kingdoms)
	GetPhysics() WeaponPhysics
	SetPhysics(physics WeaponPhysics)

	// Проризание (насколько глубоко оружие входит в объект)
	GetPenetration() float64
	SetPenetration(penetration float64)
}

// ============================================================
// WeaponPhysics - физические параметры оружия
// ============================================================

type WeaponPhysics struct {
	Weight              float64 // Вес (влияет на инерцию)
	Friction            float64 // Трение (скольжение)
	ImpactForce         float64 // Сила удара (отдача)
	MaxAngularVelocity  float64 // Максимальная скорость вращения
	AngularAcceleration float64 // Ускорение при вращении
	Damping             float64 // Демпфирование (замедление вращения)
	SmoothingPos        float64 // сила притяжаения координт оружия к целивой
	SmoothingReturn     float64 // скорость возврата угла к цели (0.05-0.2)
}
