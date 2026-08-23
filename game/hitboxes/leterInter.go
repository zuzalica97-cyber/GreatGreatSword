package hitboxes

// ============================================================
// LetterSender - объект, который может отправлять эффекты
// ============================================================

type LetterSender interface {
	// GetEffectsForTransfer - возвращает эффекты для передачи
	// Если эффектов нет или кулдаун активен — возвращает nil
	GetEffectsForTransfer(object any) []Effect

	// CanSendEffects - можно ли отправить эффекты сейчас
	CanSendEffects(object any) bool
}

// ============================================================
// LetterReceiver - объект, который может получать эффекты
// ============================================================

type LetterReceiver interface {
	// OnCollisionEffects - вызывается при получении эффектов
	// Получатель применяет эффекты к себе
	OnCollision(effects []Effect)
}
