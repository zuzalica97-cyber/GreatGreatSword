package player

import (
	"fmt"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
)

// ============================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА EffectUser
// ============================================================

func (p *PlayerLeg) AddEffect(effect hitboxes.Effect) {
	// Проверяем, есть ли уже такой эффект
	for _, e := range p.Effects {
		if e.GetType() == effect.GetType() && e.IsActive() && e.CanExtend() {
			// Суммируем время
			e.SetDuration(e.GetDuration() + effect.GetDuration())
			e.SetTimer(e.GetDuration())
			return
		}
	}
	// Если нет — добавляем новый
	p.Effects = append(p.Effects, effect)
	effect.Apply(p)
}

func (p *PlayerLeg) RemoveEffect(effectID string) {
	for i, e := range p.Effects {
		if e.GetID() == effectID {
			p.Effects = append(p.Effects[:i], p.Effects[i+1:]...)
			return
		}
	}
}

// UpdateEffects - обновляет все эффекты врага
// Возвращает true, если есть активные эффекты
func (p *PlayerLeg) UpdateEffects(dt float64) bool {
	hasActive := false
	for i := 0; i < len(p.Effects); i++ {
		effect := p.Effects[i]
		if effect.Update(p, dt) || !effect.IsActive() {
			// Эффект завершён — удаляем
			p.Effects = append(p.Effects[:i], p.Effects[i+1:]...)
			i--
		} else {
			hasActive = true
		}
	}
	return hasActive
}

func (p *PlayerLeg) GetEffects() []hitboxes.Effect {
	return p.Effects
}

func (p *PlayerLeg) HasEffect(effectType string) bool {
	for _, e := range p.Effects {
		if e.GetType() == effectType && e.IsActive() {
			return true
		}
	}
	return false
}

func (p *PlayerLeg) ClearEffects() {
	p.Effects = make([]hitboxes.Effect, 0)
}

func (p *PlayerLeg) TakeDamage(damage float64) {
	common.PlayerHelth -= damage
	if common.PlayerHelth < 0 {
		common.PlayerHelth = 0
	}
}

func (p *PlayerLeg) GetHealth() float64 {
	return common.PlayerHelth
}

func (p *PlayerLeg) GetSpeed() float64 {
	return p.CurrentMaxSpeed
}

func (p *PlayerLeg) SetSpeed(speed float64) {
	p.CurrentMaxSpeed = speed
}

func (p *PlayerLeg) SetSize(size int) {
	return
}

func (p *PlayerLeg) GetSize() int {
	return common.PlayerSize
}

func (p *PlayerLeg) SetHealth(health float64) {
	common.PlayerHelth = health
	if common.PlayerHelth < 0 {
		common.PlayerHelth = 0
	}
}

// ============================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА LetterReceiver
// ============================================================

func (p *PlayerLeg) OnCollision(effects []hitboxes.Effect) {
	for _, effect := range effects {
		fmt.Println(1)
		p.AddEffect(effect)
	}
}

// ============================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА LetterSender
// ============================================================

func (p *PlayerLeg) GetEffectsForTransfer() []hitboxes.Effect {
	for _, letter := range p.Letters {
		if !letter.CanDeliver() {
			continue
		}

		letter.Deliver()

		var clones []hitboxes.Effect
		// Возвращаем клоны эффектов для передачи
		for _, effect := range p.Effects {
			clones = append(clones, effect.Clone())
		}
		return clones
	}
	return nil
}

func (p *PlayerLeg) CanSendEffects() bool {
	// Игрок может отправлять эффекты, если они есть
	return len(p.Letters) > 0
}

func (p *PlayerLeg) OnEffectsSent() {
}
