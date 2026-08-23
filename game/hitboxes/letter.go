package hitboxes

import (
	"fmt"
	"reflect"
)

const MinCoolDownLetter = 0.1

type LetergDelivers struct {
	Target        string
	CollDownTimer float64
}

func NewDelivers(target string, coolDown float64) *LetergDelivers {
	return &LetergDelivers{
		Target:        target,
		CollDownTimer: coolDown,
	}
}

type Letter struct {
	Effects        []Effect
	MaxDeliveries  int // -1 = бесконечно
	DeliveredCount int
	IsInfinite     bool
	CoolDown       float64
	DeliversMatrix []*LetergDelivers
	Targets        []reflect.Type
}

func NewLetter(infinity bool, coolDown float64, effects []Effect, targets ...reflect.Type) *Letter {
	// Если кулдаун меньше минимального — устанавливаем минимальный
	if coolDown < MinCoolDownLetter {
		coolDown = MinCoolDownLetter
	}

	maxDeliveries := 1
	if infinity {
		maxDeliveries = -1
	}

	return &Letter{
		Effects:        effects,
		MaxDeliveries:  maxDeliveries,
		DeliveredCount: 0,
		IsInfinite:     infinity,
		CoolDown:       coolDown,
		Targets:        targets,
	}
}

// CanDeliver - можно ли отправить письмо
func (l *Letter) CanDeliver(target any) bool {
	// Проверяем, не истекло ли письмо
	if !l.IsInfinite && l.DeliveredCount >= l.MaxDeliveries {
		return false
	}

	hitbox, ok := target.(HitBoxer)

	if !ok {
		return false
	}

	for _, deliversTarget := range l.DeliversMatrix {

		if deliversTarget.Target == hitbox.GetHitBoxID() {
			return false
		}
		fmt.Println(deliversTarget.Target, hitbox.GetHitBoxID())
	}
	// Проверяем кулдаун (таймер должен быть <= 0)
	return true
}

// Deliver - отправить письмо
func (l *Letter) Deliver(target any) {
	if !l.IsInfinite {
		l.DeliveredCount++
	}

	hitbox, ok := target.(HitBoxer)

	if !ok {
		return
	}
	// Запускаем кулдаун после отправки
	l.DeliversMatrix = append(l.DeliversMatrix, NewDelivers(hitbox.GetHitBoxID(), l.CoolDown))
}

// UpdateCoolDown - обновляет таймер кулдауна
// Возвращает true, если кулдаун прошёл (таймер <= 0)
func (l *Letter) UpdateCoolDown(dt float64) bool {
	if len(l.DeliversMatrix) == 0 {
		return true
	}
	for i := 0; i < len(l.DeliversMatrix); i++ {
		l.DeliversMatrix[i].CollDownTimer -= dt

		if l.DeliversMatrix[i].CollDownTimer <= 0 {
			l.DeliversMatrix[i] = nil
			l.DeliversMatrix = append(l.DeliversMatrix[:i], l.DeliversMatrix[i+1:]...)
			i--
		}
	}
	return len(l.DeliversMatrix) == 0
}

// IsExpired - проверяет, истекло ли письмо (больше нельзя отправлять)
func (l *Letter) IsExpired() bool {
	return !l.IsInfinite && l.DeliveredCount >= l.MaxDeliveries
}

// ResetCooldown - сбрасывает кулдаун (принудительно) НЕ РЕАИЛИЗОВАННО
func (l *Letter) ResetCooldown() {
	return
}

// WhiteListLetters - проверяет, соответствует ли объект хотя бы одному интерфейсу из белого списка
// Использует reflect для универсальной проверки любых интерфейсов
func (l *Letter) WhiteListLetters(target interface{}) bool {
	if len(l.Targets) == 0 {
		return true
	}

	targetType := reflect.TypeOf(target)
	if targetType == nil {
		return false
	}

	for _, targetInterface := range l.Targets {
		if targetInterface == nil {
			continue
		}
		if targetType.Implements(targetInterface) {
			return true
		}
	}
	return false
}
