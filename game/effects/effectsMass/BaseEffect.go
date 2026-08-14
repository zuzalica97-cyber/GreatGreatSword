package effectsmass

// ============================================================
// БАЗОВАЯ СТРУКТУРА ЭФФЕКТА
// ============================================================

type BaseEffect struct {
	// ===== ИДЕНТИФИКАЦИЯ =====
	id         string // уникальный ID эффекта (генерируется при создании)
	effectType string // тип эффекта: "burn", "slow", "freeze" и т.д.

	// ===== СОСТОЯНИЕ =====
	active bool // активен ли эффект в данный момент

	// ===== ВРЕМЯ =====
	duration float64 // полная длительность эффекта (в секундах)
	timer    float64 // текущий таймер (уменьшается каждый кадр)

	// ===== ПЕРЕДАЧА ДРУГИМ ОБЪЕКТАМ =====
	maxTransfers       int // максимальное количество передач (0 = нельзя передавать)
	remainingTransfers int // сколько осталось передач

	// ===== ПОВЕДЕНИЕ =====
	canStack  bool // можно ли накладывать несколько эффектов одного типа
	canExtend bool // можно ли продлить время жизни эффекта
}

// NewBaseEffect - создаёт базовый эффект с указанными параметрами
//   - effectType: тип эффекта ("burn", "slow", "freeze")
//   - duration: длительность в секундах
//   - maxTransfers: сколько раз можно передать (0 = нельзя)
//   - canStack: можно ли накладывать несколько раз
//   - canExtend: можно ли продлить время при повторном наложении
func NewBaseEffect(effectType string, duration float64, maxTransfers int, canStack, canExtend bool) BaseEffect {
	return BaseEffect{
		id:                 effectType + "_" + string(7),
		effectType:         effectType,
		active:             true,
		duration:           duration,
		timer:              duration,
		maxTransfers:       maxTransfers,
		remainingTransfers: maxTransfers,
		canStack:           canStack,
		canExtend:          canExtend,
	}
}

// ============================================================
// МЕТОДЫ ИДЕНТИФИКАЦИИ
// ============================================================

// GetID - возвращает уникальный ID эффекта.
// Используется для поиска и удаления конкретного эффекта.
func (b *BaseEffect) GetID() string {
	return b.id
}

// GetType - возвращает тип эффекта ("burn", "slow", "freeze").
// Используется для проверки, есть ли уже такой эффект на объекте.
func (b *BaseEffect) GetType() string {
	return b.effectType
}

// ============================================================
// МЕТОДЫ СОСТОЯНИЯ
// ============================================================

// IsActive - возвращает true, если эффект активен.
// Неактивные эффекты удаляются из менеджера.
func (b *BaseEffect) IsActive() bool {
	return b.active
}

// SetActive - устанавливает состояние эффекта.
// Используется при завершении эффекта или его принудительном снятии.
func (b *BaseEffect) SetActive(active bool) {
	b.active = active
}

// ============================================================
// МЕТОДЫ ВРЕМЕНИ
// ============================================================

// GetDuration - возвращает полную длительность эффекта в секундах.
// Используется для отображения в UI или при продлении эффекта.
func (b *BaseEffect) GetDuration() float64 {
	return b.duration
}

// SetDuration - устанавливает новую длительность эффекта.
// Используется при продлении эффекта (если canExtend == true).
func (b *BaseEffect) SetDuration(duration float64) {
	b.duration = duration
}

// GetTimer - возвращает текущее значение таймера (сколько осталось секунд).
// Используется для обновления UI и проверки завершения эффекта.
func (b *BaseEffect) GetTimer() float64 {
	return b.timer
}

// SetTimer - устанавливает новое значение таймера.
// Используется при продлении эффекта или его обновлении.
func (b *BaseEffect) SetTimer(timer float64) {
	b.timer = timer
}

// ============================================================
// МЕТОДЫ ПЕРЕДАЧИ
// ============================================================

// GetMaxTransfers - возвращает максимальное количество передач.
// 0 означает, что эффект нельзя передавать другим объектам.
func (b *BaseEffect) GetMaxTransfers() int {
	return b.maxTransfers
}

// GetRemainingTransfers - возвращает оставшееся количество передач.
// Используется для проверки, можно ли ещё передать эффект.
func (b *BaseEffect) GetRemainingTransfers() int {
	return b.remainingTransfers
}

// SetRemainingTransfers - устанавливает оставшееся количество передач.
// Используется при успешной передаче эффекта другому объекту.
func (b *BaseEffect) SetRemainingTransfers(count int) {
	b.remainingTransfers = count
}

// CanTransfer - проверяет, можно ли передать эффект другому объекту.
// Возвращает true, если осталось хотя бы одна передача.
func (b *BaseEffect) CanTransfer() bool {
	return b.remainingTransfers > 0
}

// ============================================================
// МЕТОДЫ ПОВЕДЕНИЯ
// ============================================================

// CanStack - возвращает true, если можно накладывать несколько эффектов одного типа.
// Например: кровотечение можно накладывать несколько раз (стаки).
// Заморозку — нельзя (только один эффект).
func (b *BaseEffect) CanStack() bool {
	return b.canStack
}

// CanExtend - возвращает true, если можно продлить время жизни эффекта.
// Если true: при повторном наложении время складывается.
// Если false: время перезаписывается или игнорируется.
func (b *BaseEffect) CanExtend() bool {
	return b.canExtend
}

// ============================================================
// БАЗОВОЕ ОБНОВЛЕНИЕ
// ============================================================

// UpdateBase - базовое обновление эффекта.
// Уменьшает таймер на dt. Если таймер становится <= 0,
// эффект помечается как неактивный и возвращает true.
// Используется в методе Update() конкретных эффектов.
func (b *BaseEffect) UpdateBase(dt float64) bool {
	if !b.active {
		return true
	}

	b.timer -= dt

	if b.timer <= 0 {
		b.active = false
		return true
	}

	return false
}
