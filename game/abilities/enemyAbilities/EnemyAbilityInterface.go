package enemyabilities

import (
	"great-sword/game"
	"great-sword/game/hitboxes"
)

// EnemyAbility - интерфейс способности врага
type EnemyAbility interface {
	Name() string
	Update(enemy EnemyUser, dt float64, manager *hitboxes.CollisionManager) bool
	Activate(enemy EnemyUser, world game.WorldView) bool
	IsActive() bool
	GetCooldown() float64
	SetCooldown(value float64)
}

// EnemyUser - интерфейс, который должен реализовать враг для способностей
type EnemyUser interface {
	GetPosition() (float64, float64)
	SetPosition(x, y float64)
	GetSpeed() float64
	SetSpeed(speed float64)
	GetTarget() (float64, float64) // позиция игрока
	SetTarget(x, y float64)
	IsMoving() bool
	ApplyPush(x, y float64)
	GetDirection() (float64, float64)
	SetDirection(x, y float64)
	UpdateDirection(targetX, targetY float64)
	GetTargetDirection(targetX, targetY float64) (float64, float64)
}
