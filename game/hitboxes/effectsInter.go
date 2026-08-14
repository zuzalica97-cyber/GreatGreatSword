package hitboxes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

// ============================================================
// ИНТЕРФЕЙС ЭФФЕКТА
// ============================================================

// Effect - интерфейс эффекта
type Effect interface {
	// Базовые методы
	GetID() string
	GetType() string
	IsActive() bool
	SetActive(active bool)

	// Время жизни
	GetDuration() float64
	SetDuration(duration float64)
	GetTimer() float64
	SetTimer(timer float64)

	// Передача другим объектам
	GetMaxTransfers() int       // максимальное количество передач
	GetRemainingTransfers() int // осталось передач
	SetRemainingTransfers(count int)
	CanTransfer() bool // можно ли передать дальше

	// Обновление
	Update(target EffectUser, dt float64) bool // возвращает true если эффект завершён

	Clone() Effect

	// Применение
	Apply(target EffectUser)
	OnTransfer(newTarget EffectUser) // вызывается при передаче другому объекту

	// CanStack - можно ли иметь несколько эффектов одного типа
	// true = можно накладывать несколько раз (например, кровотечение)
	// false = только один экземпляр (например, заморозка)
	CanStack() bool

	// CanExtend - можно ли продлить время жизни эффекта
	// true = время складывается при повторном наложении
	// false = время перезаписывается или игнорируется
	CanExtend() bool

	Draw(screen *ebiten.Image, camera *kamera.Camera, target EffectUser)
}

// ============================================================
// ИНТЕРФЕЙС ПОЛЬЗОВАТЕЛЯ ЭФФЕКТОВ
// ============================================================

// EffectUser - объект, который может получать эффекты
type EffectUser interface {
	HitBoxer

	// Управление эффектами
	AddEffect(effect Effect)
	RemoveEffect(effectID string)
	GetEffects() []Effect
	HasEffect(effectType string) bool
	ClearEffects()

	// Получение данных для эффектов
	GetPosition() (float64, float64)
	GetSpeed() float64
	SetSpeed(speed float64)
	GetMaxSpeed() float64
	TakeDamage(damage float64)
	GetHealth() float64
	SetHealth(float64)
	GetSize() int
	SetSize(size int)
}
