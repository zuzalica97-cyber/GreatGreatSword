package game

import (
	"great-sword/game/hitboxes"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

type EntityManager struct {
	// Все сущности
	entities []Entity

	// Враги (отдельный слайс для быстрого доступа)
	enemies []Enemy

	spawners map[string]func(x, y float64) Enemy // ← реестр функций

	// Игрок
	player Player
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities: make([]Entity, 0),
		enemies:  make([]Enemy, 0),

		spawners: make(map[string]func(x float64, y float64) Enemy),
	}
}

// ============================================================
// ДОБАВЛЕНИЕ
// ============================================================

// RegisterEnemyType - регистрирует тип врага (вызывается из main)
func (em *EntityManager) RegisterEnemyType(name string, spawner func(x, y float64) Enemy) {
	em.spawners[name] = spawner
}

// SpawnEnemy - создаёт врага по типу
func (em *EntityManager) SpawnEnemy(enemyType string, x, y float64) {
	spawner, ok := em.spawners[enemyType]
	if !ok {
		log.Printf("Неизвестный тип врага: %s", enemyType)
		return
	}

	enemy := spawner(x, y)
	if enemy == nil {
		return
	}

	em.enemies = append(em.enemies, enemy)
	em.entities = append(em.entities, enemy)
}

// AddPlayer - добавляет игрока
func (em *EntityManager) AddPlayer(player Player) {
	em.player = player
	em.entities = append(em.entities, player)

}

// SpawnFromRoom - спавнит всех врагов из карточки комнаты
func (em *EntityManager) SpawnFromRoom(room *RoomCard) {
	for _, spawn := range room.Enemies {
		em.SpawnEnemy(spawn.EnemyType, spawn.X, spawn.Y)
	}
}

// ============================================================
// ОБНОВЛЕНИЕ
// ============================================================

func (em *EntityManager) Update(world WorldView, manager *hitboxes.CollisionManager) bool {
	// 1. Обновляем игрока
	if em.player != nil {
		em.player.Update(world, manager)
	}

	// 2. Обновляем врагов
	for i := 0; i < len(em.enemies); i++ {
		enemy := em.enemies[i]

		// Если враг мёртв — удаляем его из всех систем
		if !enemy.IsActive() || enemy.GetHealth() <= 0 {
			em.removeEnemy(i, manager)
			i--
			continue
		}

		enemy.Update(world, manager)
	}

	return false
}

// ============================================================
// УДАЛЕНИЕ
// ============================================================

// removeEnemy - удаляет врага из всех систем
func (em *EntityManager) removeEnemy(index int, manager *hitboxes.CollisionManager) {
	enemy := em.enemies[index]
	// 2. Вызываем OnDeath (счёт, буст и т.д.)
	enemy.OnDeath(manager)
	// 3. Удаляем из слайса врагов
	em.enemies = append(em.enemies[:index], em.enemies[index+1:]...)
}

// ============================================================
// ОТРИСОВКА
// ============================================================

func (em *EntityManager) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	for _, entity := range em.entities {

		if drawler, ok := entity.(Drawler); ok {
			drawler.Draw(screen, camera)
		}

	}
}

// ============================================================
// ПОИСК
// ============================================================

func (em *EntityManager) GetEnemies() []Enemy {
	return em.enemies
}

func (em *EntityManager) GetPlayer() Player {
	return em.player
}

func (em *EntityManager) IsActive() bool {
	return true
}

func (em *EntityManager) CollEntity() int {
	return len(em.spawners)
}

func (em *EntityManager) Tag() string {
	return "EntityManager"
}
