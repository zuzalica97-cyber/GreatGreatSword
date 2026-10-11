package enemy

import (
	"great-sword/game"
	"great-sword/game/hitboxes"
)

// ============================================================
// ФАБРИКА ВРАГОВ
// ============================================================

// EnemyCreator - функция создания врага
type EnemyCreator func(x, y float64, manager *hitboxes.CollisionManager) game.Enemy

// EnemyFactory - реестр всех типов врагов
type EnemyFactory struct {
	creators map[string]EnemyCreator
}

func NewEnemyFactory() *EnemyFactory {
	return &EnemyFactory{
		creators: make(map[string]EnemyCreator),
	}
}

// Register - регистрирует тип врага
func (f *EnemyFactory) Register(enemyType string, creator EnemyCreator) {
	f.creators[enemyType] = creator
}

// Create - создаёт врага по типу
func (f *EnemyFactory) Create(enemyType string, x, y float64, manager *hitboxes.CollisionManager) game.Enemy {
	creator, ok := f.creators[enemyType]
	if !ok {
		return nil
	}
	return creator(x, y, manager)
}

// Has - проверяет, зарегистрирован ли тип
func (f *EnemyFactory) Has(enemyType string) bool {
	_, ok := f.creators[enemyType]
	return ok
}

// RegisterAll - регистрирует всех врагов игры
func (f *EnemyFactory) RegisterAll() {
	f.Register("berserk", func(x, y float64, manager *hitboxes.CollisionManager) game.Enemy {
		return SpawnBreserk(x, y, manager)
	})
}
