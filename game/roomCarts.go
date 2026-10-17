package game

// ============================================================
// КАРТОЧКА КОМНАТЫ
// ============================================================
// RoomCard - карточка комнаты
type RoomCard struct {
	Enemies []EnemySpawn
}

// EnemySpawn - данные о спавне одного врага
type EnemySpawn struct {
	EnemyType string  // "berserk", "pathetic", "hater"
	X, Y      float64 // координаты спавна
}

func NewRoomCard() *RoomCard {
	return &RoomCard{
		Enemies: make([]EnemySpawn, 0),
	}
}

// AddEnemy - добавляет врага в карточку
func (rc *RoomCard) AddEnemy(enemyType string, x, y float64) {
	rc.Enemies = append(rc.Enemies, EnemySpawn{
		EnemyType: enemyType,
		X:         x,
		Y:         y,
	})
}
