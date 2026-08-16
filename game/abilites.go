package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Ability interface {
	// Название способности
	Name() string

	// Обновление каждый кадр (возвращает true, если способность активна)
	Update(world WorldView) bool

	// Отрисовка (опционально, для UI)
	Draw(screen *ebiten.Image)

	Activate(world WorldView)

	// Проверка, активна ли способность
	IsActive() bool

	// Время до повторного использования (кулдаун)
	CooldownAbilites() float64
	SetCooldown(float64)
}

type PlayerLegInter interface {
	ApplyForce(vx, vy float64)
	GetDirection() (float64, float64)
	GetPosition() (float64, float64)
	SetPosition(x, y float64)
	GetSpeedXY() (float64, float64)
	SetSpeedXY(vx, vy float64)
	GetMaxSpeed() float64
	SetMaxSpeed(speed float64)
	IsMoving() bool
	GetHealth() float64
	SetHealth(float64)
}

type PlayerHeadInter interface {
	GetRotationSpeed() float64
	SetRotationSpeed(speed float64)
	GetAngle() float64
	SetAngle(angle float64)
	GetPosition() (float64, float64)
	GetCurrentRotSpeed() float64
	SetCurrentRotSpeed(speed float64)
	GetAceleration() float64
	SetAceleration(axel float64)
	GetDeAceleration() float64
	SetDeAceleration(dexel float64)
}
