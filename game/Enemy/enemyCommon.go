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

func PatheticCooldown(enemy *OnePath) {
	enemy.PathericCooldownActive = true
	enemy.PathericCooldownTimer = PatheticCooldownTime
}

func HatersCooldown(enemy *Haits) {
	enemy.HatersCooldownActive = true
	enemy.HatersCooldownTimer = HatersCooldownTime
}

func CheckCollisionWithPlayer(enemyX, enemyY, enemySize float64, world game.WorldView) bool {

	var playerX float64
	var playerY float64

	for _, entity := range world.SearchEntities("playerLeg") {
		leg := entity.(*player.PlayerLeg)

		playerX = leg.Position.Px
		playerY = leg.Position.Py

	}
	return enemyX < playerX+common.PlayerSize &&
		enemyX+enemySize > playerX &&
		enemyY < playerY+common.PlayerSize &&
		enemyY+enemySize > playerY
}

func CheckCollisionWithAttachment(enemyX, enemyY, enemySize float64, attX, attY, attW, attH, attAngle float64) bool {
	cos := math.Cos(attAngle)
	sin := math.Sin(attAngle)

	enemyCenterX := enemyX + enemySize/2
	enemyCenterY := enemyY + enemySize/2

	localX := (enemyCenterX-attX)*cos + (enemyCenterY-attY)*sin
	localY := -(enemyCenterX-attX)*sin + (enemyCenterY-attY)*cos

	halfW := attW / 2
	halfH := attH / 2

	return math.Abs(localX) < halfW+enemySize/2 && math.Abs(localY) < halfH+enemySize/2
}

func RangomSpawnInWall() (float64, float64) {
	side := rand.IntN(4)

	var x, y float64

	switch side {
	case 0: //левый крвй
		x = -float64(common.PatheticSize)
		y = float64(rand.IntN(common.ScreenHeight))
	case 1: //левый крвй
		x = float64(common.ScreenWidth)
		y = float64(rand.IntN(common.ScreenHeight))
	case 2: //левый крвй
		x = float64(rand.IntN(common.ScreenWidth))
		y = -float64(common.PatheticSize)
	case 3: //левый крвй
		x = float64(rand.IntN(common.ScreenWidth))
		y = float64(common.ScreenHeight)
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
