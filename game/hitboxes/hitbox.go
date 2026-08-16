package hitboxes

import (
	"github.com/setanarut/coll"
	v "github.com/setanarut/v"
)

// Фиксированные значения
const auraRadius = 30.0
const auraForce = 3.0

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

	// GetWeight - возвращает вес объекта
	// Чем больше вес, тем сложнее его сдвинуть
	// Вес 1 — лёгкий, вес 10 — тяжёлый
	GetWeight() float64

	// GetDensity - возвращает плотность объекта
	// 0.0 - полностью проницаемый (как призрак)
	// 0.5 - частично проницаемый (можно войти с трудом)
	// 1.0 - полностью твёрдый (как стена, нельзя войти)
	GetDensity() float64

	// HasAura - имеет ли объект ауру отталкивания
	HasAura() bool

	// AffectedByAura - реагирует ли объект на чужую ауру
	AffectedByAura() bool
}

// ============================================================
// ВРАЩАЮСЕИСЯ ОБЬЕКТЫ
// ============================================================

type RotatableHitBoxer interface {
	HitBoxer
	GetOBB() (centerX, centerY, halfW, halfH, angle float64) // центр, половинные размеры, угол
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
			cm.objects[i] = nil
			cm.objects = append(cm.objects[:i], cm.objects[i+1:]...)
			break
		}
	}
}

// Update - проверяет все коллизии между объектами
func (cm *CollisionManager) Update() {
	for i := 0; i < len(cm.objects); i++ {
		for j := i + 1; j < len(cm.objects); j++ {
			obj1 := cm.objects[i]
			obj2 := cm.objects[j]

			if !obj1.IsActive() || !obj2.IsActive() {
				continue
			}

			cm.checkCollision(obj1, obj2)

			// ===== АУРА ОТТАЛКИВАНИЯ (только если есть аура) =====
			if obj1.HasAura() || obj2.HasAura() {
				pushX1, pushY1, pushX2, pushY2 := CalculateAuraPush(obj1, obj2)

				if !obj1.IsStatic() && obj1.AffectedByAura() {
					obj1.ApplyPush(pushX1, pushY1)
				}
				if !obj2.IsStatic() && obj2.AffectedByAura() {
					obj2.ApplyPush(pushX2, pushY2)
				}
			}

		}
	}

}

// checkCollision - определяет тип коллизии и вызывает нужную проверку
func (cm *CollisionManager) checkCollision(obj1, obj2 HitBoxer) {

	//проверка евляютса ли они вращяюсемися
	rot1, isRot1 := obj1.(RotatableHitBoxer)
	rot2, isRot2 := obj2.(RotatableHitBoxer)

	switch {
	case isRot1 && isRot2:
		// TODO: реализовать OBB vs OBB, если нужно
		// cm.checkOBBvsOBB(rot1, rot2)
	case isRot1:
		cm.checkOBBvsAABB(rot1, obj2)
	case isRot2:
		cm.checkOBBvsAABB(rot2, obj1)
	default:
		cm.checkAABBvsAABB(obj1, obj2)
	}
}

// checkOBBvsOBB - проверка столкновения двух вращающихся объектов
// ТУТ ДОЛЖНО БЫТЬ ПРОВЕРКА ДВУХ ВРОЩЯЮЩИХСЯ ОБЬЕКТВОВ

// checkOBBvsAABB - проверка столкновения вращающегося и статичного объекта
func (cm *CollisionManager) checkOBBvsAABB(rot RotatableHitBoxer, stat HitBoxer) {
	cx, cy, hw, hh, angle := rot.GetOBB()

	obb := &coll.OBB{
		Pos:   v.Vec{X: cx, Y: cy},
		Half:  v.Vec{X: hw, Y: hh},
		Angle: angle,
	}

	px, py, hw2, hh2 := stat.GetAABB()
	aabb := &coll.AABB{
		Pos:  v.Vec{X: px, Y: py},
		Half: v.Vec{X: hw2, Y: hh2},
	}

	hit := &coll.Hit{}

	// 3. Используем готовую функцию из библиотеки coll
	// Она вернёт true, если AABB и OBB пересекаются
	if coll.BoxOrientedBoxOverlap(aabb, obb) {
		// Если столкновение есть, обрабатываем его
		// Так как у нас нет hit-информации, создаём пустой хит или
		// вызываем resolveCollision без hit-данных
		cm.resolveCollision(rot, stat, hit)
	}
}

// checkAABBvsAABB - проверка столкновения двух статичных объектов
func (cm *CollisionManager) checkAABBvsAABB(obj1, obj2 HitBoxer) {
	px1, py1, hw1, hh1 := obj1.GetAABB()
	px2, py2, hw2, hh2 := obj2.GetAABB()

	aabb1 := &coll.AABB{
		Pos:  v.Vec{X: px1, Y: py1},
		Half: v.Vec{X: hw1, Y: hh1},
	}
	aabb2 := &coll.AABB{
		Pos:  v.Vec{X: px2, Y: py2},
		Half: v.Vec{X: hw2, Y: hh2},
	}

	hit := &coll.Hit{}
	if coll.BoxBoxOverlap(aabb1, aabb2, hit) {
		cm.resolveCollision(obj1, obj2, hit)
	}
}

// resolveCollision - обрабатывает столкновение между двумя объектами
// В hitboxes/collision_manager.go

func (cm *CollisionManager) resolveCollision(a, b HitBoxer, hit *coll.Hit) {
	if hit == nil {
		return
	}

	normal := hit.Normal
	penetration := hit.Data

	ABSender(a, b)
	ABSender(b, a)

	if a.IsStatic() && b.IsStatic() {
		return
	}

	// ===== НОВАЯ ЛОГИКА С УЧЁТОМ ВЕСА И ПЛОТНОСТИ =====

	weightA := a.GetWeight()
	weightB := b.GetWeight()
	densityA := a.GetDensity()
	densityB := b.GetDensity()

	// Общий вес для расчёта отталкивания
	totalWeight := weightA + weightB

	// Коэффициенты отталкивания (чем больше вес, тем меньше отталкивание)
	pushFactorA := weightB / totalWeight
	pushFactorB := weightA / totalWeight

	// Коэффициент плотности (чем выше плотность, тем сильнее отталкивание)
	densityFactor := (densityA + densityB) / 2

	// Если плотность очень низкая (< 0.3) — почти не отталкиваем
	if densityFactor < 0.3 {
		return // объекты проходят друг через друга
	}

	// Базовая сила отталкивания
	pushStrength := 0.3 * densityFactor

	// Применяем отталкивание с учётом веса и плотности
	if !a.IsStatic() {
		pushX := -normal.X * penetration * pushStrength * pushFactorA
		pushY := -normal.Y * penetration * pushStrength * pushFactorA
		a.ApplyPush(pushX, pushY)
	}

	if !b.IsStatic() {
		pushX := normal.X * penetration * pushStrength * pushFactorB
		pushY := normal.Y * penetration * pushStrength * pushFactorB
		b.ApplyPush(pushX, pushY)
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

// ============================================================
// CompactSlice - удаляет все nil элементы из слайса указателей
// ============================================================
//
// [T any]     - "Дженерик": T может быть ЛЮБЫМ типом (заполнитель)
// []*T        - слайс УКАЗАТЕЛЕЙ на тип T (работаем с оригиналами)
// []*T        - возвращаем такой же слайс УКАЗАТЕЛЕЙ
//
// Зачем:
// 1. Удаляем мёртвые объекты (enemy = nil)
// 2. Освобождаем память
// 3. Избавляемся от мусора в слайсе
// ============================================================

func CompactSlice[T any](slice []*T) []*T {
	// Создаём новый слайс с ТАКОЙ ЖЕ вместимостью (cap),
	// чтобы избежать лишних аллокаций
	compacted := make([]*T, 0, len(slice))

	// Проходим по всем элементам
	for _, item := range slice {
		// Если элемент НЕ nil - сохраняем
		if item != nil {
			compacted = append(compacted, item)
		}
		// Если nil - пропускаем (удаляем)
	}

	// Возвращаем чистый слайс без nil
	return compacted
}

// В hitboxes/collision_manager.go

// Cleanup - удаляет все nil объекты из менеджера
func (cm *CollisionManager) Cleanup() {
	compacted := make([]HitBoxer, 0, len(cm.objects))
	for _, obj := range cm.objects {
		if obj != nil {
			compacted = append(compacted, obj)
		}
	}
	cm.objects = compacted
}

func (cm *CollisionManager) LenColisionObjects() int {
	return len(cm.objects)
}
