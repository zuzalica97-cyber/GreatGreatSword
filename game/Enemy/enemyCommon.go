package enemy

import (
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/player"
	swords "great-sword/game/player/Swords"
	"math"
	"math/rand/v2"
)

var PatheticCooldownTime = 1.0
var HatersCooldownTime = 1.0

func getPlayerPosition(worldView game.WorldView) (float64, float64) {
	for _, pleg := range worldView.SearchEntities("playerLeg") {
		p := pleg.(game.PlayerLegInter)
		return p.GetPosition()

	}
	return 0, 0
}

func RangomSpawnInWall(size int) (float64, float64) {
	side := rand.IntN(4)

	var x, y float64

	switch side {
	case 0: //левый крвй
		x = -float64(size)
		y = float64(rand.IntN(int(common.RoomHeight)))
	case 1: //левый крвй
		x = float64(common.RoomWidth)
		y = float64(rand.IntN(int(common.RoomHeight)))
	case 2: //левый крвй
		x = float64(rand.IntN(int(common.RoomWidth)))
		y = -float64(size)
	case 3: //левый крвй
		x = float64(rand.IntN(int(common.RoomWidth)))
		y = float64(common.RoomHeight)
	}
	return x, y
}

func WhereThePlayer(world game.WorldView) (float64, float64, float64, float64, float64, float64, float64) {

	var playerX, playerY, attahmentX, attahmentY, attahmentW, attahmentH, attahmentAngle float64

	for _, entity := range world.SearchEntities("playerLeg") {
		pLeg := entity.(*player.PlayerLeg)

		playerX = pLeg.Position.Px + common.PlayerSize/2
		playerY = pLeg.Position.Py + common.PlayerSize/2

		for _, entity := range world.SearchEntities("blueSword") {
			sword := entity.(*swords.BlueSword)

			attahmentX = sword.Position.Px
			attahmentY = sword.Position.Py
			attahmentW = common.SwordAttachmentWidth
			attahmentH = common.SwordAttachmentHeight
			attahmentAngle = sword.Angle * math.Pi / 180
		}
	}
	return playerX, playerY, attahmentX, attahmentY, attahmentW, attahmentH, attahmentAngle
}

// MoveEnemyTowardsPlayer - передвигает врага в сторону игрока
// Параметры:
//   - enemyX, enemyY: текущая позиция врага
//   - playerX, playerY: текущая позиция игрока
//   - speed: скорость передвижения врага
//   - dt: дельта времени
//
// Возвращает: новые координаты врага (newX, newY)
func MoveEnemyTowardsPlayer(enemyX, enemyY, playerX, playerY, speed, dt float64) (float64, float64) {
	// Вычисляем направление к игроку
	dx := playerX - enemyX
	dy := playerY - enemyY
	distance := math.Sqrt(dx*dx + dy*dy)

	// Если враг уже на месте или скорость 0 — не двигаемся
	if distance <= 0.01 || speed == 0 {
		return enemyX, enemyY
	}

	// Нормализуем направление и применяем скорость
	dirX := dx / distance
	dirY := dy / distance

	newX := enemyX + dirX*speed*dt
	newY := enemyY + dirY*speed*dt

	return newX, newY
}
