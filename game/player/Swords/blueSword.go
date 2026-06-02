package swords

import (
	"fmt"
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/player"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Entity = (*BlueSword)(nil)

type BlueSword struct {
	Position         common.PointPlayer
	TargetX, TargetY float64
	Angle            float64
	TargetAngle      float64
	SnowCollision    bool
}

func NewBlueSword(world game.WorldView) *BlueSword {
	b := &BlueSword{}

	b.UpdateAttachmentTarget(world)
	b.Position.Px = b.TargetX
	b.Position.Py = b.TargetY
	b.Angle = b.TargetAngle

	return b
}

func (b *BlueSword) UpdateAttachmentWithDelay() {
	b.Position.Px = b.Position.Px*0 +
		b.TargetX*(1-0)
	b.Position.Py = b.Position.Py*0 +
		b.TargetY*(1-0)

	angleDiff := b.TargetAngle - b.Angle

	if angleDiff > 180 {
		angleDiff -= 360
	} else if angleDiff < -180 {
		angleDiff += 360
	}

	b.Angle += angleDiff * (1 - common.SwordAttachmentSmoothing)

	b.Angle = math.Mod(b.Angle, 360)
	if b.Angle < 0 {
		b.Angle += 360
	}
}

func (b *BlueSword) GetAttachmentCorners() [4]struct{ X, Y float64 } {
	angleRed := b.Angle * math.Pi / 180
	cos := math.Cos(angleRed)
	sin := math.Sin(angleRed)

	halfW := common.SwordAttachmentWidth / 2
	halfH := common.SwordAttachmentHeight / 2

	localCorners := [4][2]float64{
		{-float64(halfW), -float64(halfH)},
		{float64(halfW), -float64(halfH)},
		{float64(halfW), float64(halfH)},
		{-float64(halfW), float64(halfH)}, //ДОДЕЛАТЬ ОТ СЮДА
	}

	var corners [4]struct{ X, Y float64 }

	for i, local := range localCorners {

		rotatedX := local[0]*cos - local[1]*sin
		rotatedY := local[0]*sin + local[1]*cos

		corners[i].X = b.Position.Px + rotatedX
		corners[i].Y = b.Position.Py + rotatedY
	}

	return corners
}

func (b *BlueSword) CheckCollisionWithWalls() bool {
	collided := false

	corners := b.GetAttachmentCorners()

	for _, corner := range corners {
		if corner.X < 0 || corner.X > common.ScreenWidth ||
			corner.Y < 0 || corner.Y > common.ScreenHeight {
			collided = true
		}
	}
	return collided
}

func (b *BlueSword) UpdateAttachmentTarget(worldView game.WorldView) {
	for _, entity := range worldView.SearchEntities("playerLeg") {
		leg := entity.(*player.PlayerLeg)

		centrX := leg.Position.Px + common.PlayerHeadSize/2
		centrY := leg.Position.Py + common.PlayerHeadSize/2

		for _, entity := range worldView.SearchEntities("playerHead") {
			head := entity.(*player.PlayerHead)

			angleRad := head.Angle * math.Pi / 180

			distanceFromCenter := common.PlayerHeadSize/2 + common.SwordAttachmentWidth/2

			offsetX := math.Cos(angleRad) * float64(distanceFromCenter)
			offsetY := math.Sin(angleRad) * float64(distanceFromCenter)

			b.TargetX = centrX + offsetX
			b.TargetY = centrY + offsetY

			b.TargetAngle = head.Angle
		}
	}
}

func (b *BlueSword) ResolveCollision(world game.WorldView) {
	for _, entity := range world.SearchEntities("playerLeg") {
		Leg := entity.(*player.PlayerLeg)

		corners := b.GetAttachmentCorners()

		var minX, maxX, minY, maxY float64

		minX = corners[0].X
		maxX = corners[0].X
		minY = corners[0].Y
		maxY = corners[0].Y

		for _, c := range corners {
			minX = math.Min(minX, c.X)
			maxX = math.Max(maxX, c.X)
			minY = math.Min(minY, c.Y)
			maxY = math.Max(maxY, c.Y)
		}

		if ebiten.IsKeyPressed(ebiten.KeyO) {
			fmt.Println(maxX)
		}

		var Mnogitel = 15.0

		var gool = 1.0

		if player.ForwardActive {
			gool = 10.0
			Mnogitel = 0.05
		}

		if player.ForwardActive && player.BoostActive {
			gool = 15.0
			Mnogitel = 0.05
		}

		offsetX := 0.0
		offsetY := 0.0

		if minX < 0 {
			offsetX = -minX * Mnogitel
		} else if maxX > common.ScreenWidth {
			offsetX = (common.ScreenWidth - maxX) * Mnogitel
		}
		if minY < 0 {
			offsetY = -minY * Mnogitel
		} else if maxY > common.ScreenHeight {
			offsetY = (common.ScreenHeight - maxY) * Mnogitel
		}

		b.Position.Px += offsetX * 0.025
		b.Position.Py += offsetY * 0.025

		Leg.Speed.Vx += offsetX * gool
		Leg.Speed.Vy += offsetY * gool
	}

}

func (b *BlueSword) Update(worldView game.WorldView) bool {

	b.UpdateAttachmentTarget(worldView)

	b.UpdateAttachmentWithDelay()

	b.ResolveCollision(worldView)

	return false
}

func (b *BlueSword) Draw(screen *ebiten.Image) {
}

func (b *BlueSword) Tag() string {
	return "blueSword"
}
