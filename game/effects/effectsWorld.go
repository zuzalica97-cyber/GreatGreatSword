package effects

import (
	"great-sword/game/hitboxes"
	"sync"
)

// ============================================================
// МЕНЕДЖЕР ЭФФЕКТОВ
// ============================================================

type EffectManager struct {
	effects map[string][]hitboxes.Effect // objectID -> список эффектов
	mu      sync.RWMutex
}

func NewEffectManager() *EffectManager {
	return &EffectManager{
		effects: make(map[string][]hitboxes.Effect),
	}
}

// AddEffect - добавляет эффект объекту
func (em *EffectManager) AddEffect(obj hitboxes.HitBoxer, effect hitboxes.Effect) {
	if obj == nil || effect == nil {
		return
	}

	objID := obj.GetHitBoxID()
	em.mu.Lock()
	defer em.mu.Unlock()

	// Добавляем новый эффект
	em.effects[objID] = append(em.effects[objID], effect)

	// Применяем эффект
	if user, ok := obj.(hitboxes.EffectUser); ok {
		effect.Apply(user)
	}
}

// RemoveEffect - удаляет эффект у объекта
func (em *EffectManager) RemoveEffect(obj hitboxes.HitBoxer, effectID string) {
	if obj == nil {
		return
	}

	objID := obj.GetHitBoxID()
	em.mu.Lock()
	defer em.mu.Unlock()

	effects := em.effects[objID]
	for i, e := range effects {
		if e.GetID() == effectID {
			em.effects[objID] = append(effects[:i], effects[i+1:]...)
			return
		}
	}
}

// Update - обновляет все эффекты
func (em *EffectManager) Update(dt float64) {
	em.mu.Lock()
	defer em.mu.Unlock()

	for objID, effects := range em.effects {
		for i := 0; i < len(effects); i++ {
			effect := effects[i]

			// Обновляем эффект
			if user, ok := getEffectUser(objID); ok {
				if effect.Update(user, dt) {
					// Эффект завершён
					em.effects[objID] = append(effects[:i], effects[i+1:]...)
					i--
				}
			}
		}
	}
}

// GetEffects - возвращает все эффекты объекта
func (em *EffectManager) GetEffects(obj hitboxes.HitBoxer) []hitboxes.Effect {
	if obj == nil {
		return nil
	}

	objID := obj.GetHitBoxID()
	em.mu.RLock()
	defer em.mu.RUnlock()

	return em.effects[objID]
}

// Clear - удаляет все эффекты у объекта
func (em *EffectManager) Clear(obj hitboxes.HitBoxer) {
	if obj == nil {
		return
	}

	objID := obj.GetHitBoxID()
	em.mu.Lock()
	defer em.mu.Unlock()

	delete(em.effects, objID)

	if user, ok := obj.(hitboxes.EffectUser); ok {
		user.ClearEffects()
	}
}

// Helper для получения EffectUser
func getEffectUser(objID string) (hitboxes.EffectUser, bool) {
	// Здесь нужна связь с вашим EntityManager
	// Можно хранить map[id]EffectUser в EffectManager
	return nil, false
}
