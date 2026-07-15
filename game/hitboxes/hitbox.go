package hitboxes

import (
	"github.com/setanarut/coll"
	"github.com/setanarut/v"
)

// HitBoxer - интерфейс для всех сущностей, у которых есть коллизия
type HitBoxer interface {
	// GetAABB возвращает AABB для проверки коллизий
	// (центр и половинные размеры)
	GetAABB() (posX, posY, halfW, halfH float64)

	// GetHitBoxID возвращает уникальный ID для идентификации
	GetHitBoxID() string

	// IsActive проверяет, активен ли объект для коллизий
	IsActive() bool

	// IsStatic проверяет, статичен ли объект (стена, платформа)
	// Если true - объект не двигается при отталкивании
	IsStatic() bool

	// ApplyPush применяет силу отталкивания (сдвиг)
	ApplyPush(x, y float64)

	// OnCollision вызывается при столкновении с другим объектом
	OnCollision(other HitBoxer)
}

// CollisionManager - управляет коллизиями между всеми HitBoxer
type CollisionManager struct {
	objects []HitBoxer
}

func NewCollisionManager() *CollisionManager {
	return &CollisionManager{
		objects: make([]HitBoxer, 0),
	}
}

// AddObject - добавляет объект в менеджер коллизий
func (cm *CollisionManager) AddObject(obj HitBoxer) {
	cm.objects = append(cm.objects, obj)
}

// RemoveObject - удаляет объект из менеджера коллизий
func (cm *CollisionManager) RemoveObject(obj HitBoxer) {
	for i, o := range cm.objects {
		if o == obj {
			cm.objects = append(cm.objects[:i], cm.objects[i+1:]...)
			break
		}
	}
}

// Update - проверяет все коллизии между объектами
func (cm *CollisionManager) Update() {
	// Проходим по всем парам объектов
	for i := 0; i < len(cm.objects); i++ {
		for j := i + 1; j < len(cm.objects); j++ {
			obj1 := cm.objects[i]
			obj2 := cm.objects[j]

			// Пропускаем неактивные объекты
			if !obj1.IsActive() || !obj2.IsActive() {
				continue
			}

			// === ПОЛУЧАЕМ AABB ===
			px1, py1, hw1, hh1 := obj1.GetAABB()
			px2, py2, hw2, hh2 := obj2.GetAABB()

			// === СОЗДАЁМ AABB ===
			aabb1 := &coll.AABB{
				Pos:  v.Vec{X: px1, Y: py1},
				Half: v.Vec{X: hw1, Y: hh1},
			}
			aabb2 := &coll.AABB{
				Pos:  v.Vec{X: px2, Y: py2},
				Half: v.Vec{X: hw2, Y: hh2},
			}

			// === ПРОВЕРЯЕМ СТОЛКНОВЕНИЕ ===
			hit := &coll.Hit{}
			if coll.BoxBoxOverlap(aabb1, aabb2, hit) {
				// === ОБРАБАТЫВАЕМ СТОЛКНОВЕНИЕ ===
				cm.resolveCollision(obj1, obj2, hit)
			}
		}
	}
}

// resolveCollision - обрабатывает столкновение между двумя объектами
func (cm *CollisionManager) resolveCollision(a, b HitBoxer, hit *coll.Hit) {
	normal := hit.Normal
	penetration := hit.Data

	// Вызываем OnCollision для обоих объектов
	a.OnCollision(b)
	b.OnCollision(a)

	// Если оба статичны — ничего не делаем
	if a.IsStatic() && b.IsStatic() {
		return
	}

	// Отталкиваем объекты
	if !a.IsStatic() {
		a.ApplyPush(-normal.X*penetration/2, -normal.Y*penetration/2)
	}
	if !b.IsStatic() {
		b.ApplyPush(normal.X*penetration/2, normal.Y*penetration/2)
	}
}

// Clear - очищает все объекты
func (cm *CollisionManager) Clear() {
	cm.objects = make([]HitBoxer, 0)
}

// GetObjects - возвращает все объекты
func (cm *CollisionManager) GetObjects() []HitBoxer {
	return cm.objects
}
