package hitboxes

const MinCoolDownLetter = 0.1

type Letter struct {
	Effects        []Effect
	MaxDeliveries  int // -1 = бесконечно
	DeliveredCount int
	IsInfinite     bool
	CoolDown       float64
	Timer          float64
}

func NewLetter(infinity bool, coolDown float64, effects ...Effect) *Letter {
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
		Timer:          0, // таймер начинается с 0 (готов к отправке)
	}
}

// CanDeliver - можно ли отправить письмо
func (l *Letter) CanDeliver() bool {
	// Проверяем, не истекло ли письмо
	if !l.IsInfinite && l.DeliveredCount >= l.MaxDeliveries {
		return false
	}
	// Проверяем кулдаун (таймер должен быть <= 0)
	return l.Timer <= 0
}

// Deliver - отправить письмо
func (l *Letter) Deliver() {
	if !l.IsInfinite {
		l.DeliveredCount++
	}
	// Запускаем кулдаун после отправки
	l.Timer = l.CoolDown
}

// UpdateCoolDown - обновляет таймер кулдауна
// Возвращает true, если кулдаун прошёл (таймер <= 0)
func (l *Letter) UpdateCoolDown(dt float64) bool {
	if l.Timer > 0 {
		l.Timer -= dt
		if l.Timer < 0 {
			l.Timer = 0
		}
	}
	return l.Timer <= 0
}

// IsExpired - проверяет, истекло ли письмо (больше нельзя отправлять)
func (l *Letter) IsExpired() bool {
	return !l.IsInfinite && l.DeliveredCount >= l.MaxDeliveries
}

// ResetCooldown - сбрасывает кулдаун (принудительно)
func (l *Letter) ResetCooldown() {
	l.Timer = 0
}
