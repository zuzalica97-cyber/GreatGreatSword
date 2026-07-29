package game

import (
	"great-sword/game/hitboxes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

// ============================================================
// МЕНЕДЖЕР СУЩНОСТЕЙ
// ============================================================

type EntityManager struct {
	// Все сущности
	entities []Entity

	// Отдельные списки для быстрого доступа
	enemies     []Enemy
	projectiles []Projectile
	player      Player

	// Менеджер коллизий
	collisionManager *hitboxes.CollisionManager
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities:         make([]Entity, 0),
		enemies:          make([]Enemy, 0),
		projectiles:      make([]Projectile, 0),
		collisionManager: hitboxes.NewCollisionManager(),
	}
}

// ============================================================
// ДОБАВЛЕНИЕ СУЩНОСТЕЙ
// ============================================================

// AddEntity - добавляет любую сущность
func (em *EntityManager) AddEntity(entity Entity) {
	em.entities = append(em.entities, entity)

	// Если это враг — добавляем в список врагов и в коллизии
	if enemy, ok := entity.(Enemy); ok {
		em.enemies = append(em.enemies, enemy)
		em.collisionManager.AddObject(enemy)
	}

	// Если это пуля — добавляем в список пуль и в коллизии
	if proj, ok := entity.(Projectile); ok {
		em.projectiles = append(em.projectiles, proj)
		em.collisionManager.AddObject(proj)
	}

	// Если это игрок — сохраняем и добавляем в коллизии
	if player, ok := entity.(Player); ok {
		em.player = player
		em.collisionManager.AddObject(player)
	}
}

// ============================================================
// ОБНОВЛЕНИЕ
// ============================================================

// Update - обновляет все сущности и проверяет коллизии
func (em *EntityManager) Update(world WorldView, manager *hitboxes.CollisionManager) bool {
	// 1. Обновляем все сущности
	for _, entity := range em.entities {
		if !entity.IsActive() {
			continue
		}
		if entity.Update(world, manager) {
			return true // игра окончена
		}
	}

	// 2. Проверяем коллизии
	em.collisionManager.Update()

	// 3. Обрабатываем результаты коллизий
	em.processCollisions(world)

	return false
}

// ============================================================
// ОБРАБОТКА КОЛЛИЗИЙ
// ============================================================

func (em *EntityManager) processCollisions(world WorldView) {
	// Все коллизии уже обработаны через OnCollision
	// Но здесь можно добавить дополнительную логику

	// Например: проверить, не умер ли враг
	for i := 0; i < len(em.enemies); i++ {
		enemy := em.enemies[i]
		if !enemy.IsActive() {
			// Враг мёртв — вызываем OnDeath
			enemy.OnDeath(world)
		}
	}
}

// ============================================================
// ОТРИСОВКА
// ============================================================

func (em *EntityManager) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	for _, entity := range em.entities {
		if entity.IsActive() {
			entity.Draw(screen, camera)
		}
	}
}

// ============================================================
// ПОИСК
// ============================================================

func (em *EntityManager) SearchEntities(tag string) []Entity {
	result := make([]Entity, 0)
	for _, entity := range em.entities {
		if entity.Tag() == tag && entity.IsActive() {
			result = append(result, entity)
		}
	}
	return result
}

func (em *EntityManager) GetEnemies() []Enemy {
	return em.enemies
}

func (em *EntityManager) GetPlayer() Player {
	return em.player
}

// ============================================================
// ОЧИСТКА
// ============================================================
