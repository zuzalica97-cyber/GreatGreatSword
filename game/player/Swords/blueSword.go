package swords

import (
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
	"great-sword/game/player"
	"image/color"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

var _ game.Entity = (*BlueSword)(nil)

type BlueSword struct {
	Position         common.PointPlayer
	TargetX, TargetY float64
	Angle            float64
	TargetAngle      float64
	SnowCollision    bool
	Texture          *ebiten.Image
}

func NewBlueSword(world game.WorldView, manager *hitboxes.CollisionManager) *BlueSword {
	b := &BlueSword{}

	b.UpdateAttachmentTarget(world)
	b.Position.Px = b.TargetX
	b.Position.Py = b.TargetY
	b.Angle = b.TargetAngle

	if manager != nil {
		manager.AddObject(b)
	}

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
		if corner.X < 0 || corner.X > common.RoomWidth ||
			corner.Y < 0 || corner.Y > common.RoomHeight {
			collided = true
		}
	}
	return collided
}

func (b *BlueSword) UpdateAttachmentTarget(worldView game.WorldView) {
	for _, entity := range worldView.SearchEntities("playerLeg") {
		p := entity.(*player.PlayerLeg)

		centrX := p.Position.Px + common.PlayerSize
		centrY := p.Position.Py + common.PlayerSize

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

		var Mnogitel = 15.0

		var gool = 1.0

		offsetX := 0.0
		offsetY := 0.0

		if minX < 0 {
			offsetX = -minX * Mnogitel
		} else if maxX > common.RoomWidth {
			offsetX = (common.RoomWidth - maxX) * Mnogitel
		}
		if minY < 0 {
			offsetY = -minY * Mnogitel
		} else if maxY > common.RoomHeight {
			offsetY = (common.RoomHeight - maxY) * Mnogitel
		}

		b.Position.Px += offsetX * 0.025
		b.Position.Py += offsetY * 0.025

		Leg.Speed.Vx += offsetX * gool
		Leg.Speed.Vy += offsetY * gool
	}

}

func (b *BlueSword) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {

	if player.SwordIxist {

		b.UpdateAttachmentTarget(worldView)

		b.UpdateAttachmentWithDelay()

		b.ResolveCollision(worldView)
	} else {
		b.Position.Px = 3000
		b.Position.Py = 3000
	}

	return false
}

func (b *BlueSword) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	Color := color.RGBA{0, 100, 200, 255}

	if player.BoostTimer > 0 {
		Color = color.RGBA{200, 50, 50, 255}
	}

	// Создаём временное изображение для меча
	swordImg := ebiten.NewImage(int(common.SwordAttachmentWidth), int(common.SwordAttachmentHeight))
	swordImg.Fill(Color)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-common.SwordAttachmentWidth/2, -common.SwordAttachmentHeight/2)
	op.GeoM.Rotate(b.Angle * math.Pi / 180)
	op.GeoM.Translate(b.Position.Px, b.Position.Py) // мировые координаты

	// Камера сама применит смещение
	camera.Draw(swordImg, op, screen)
}

// Для AABB (приближение)
func (b *BlueSword) GetAABB() (posX, posY, halfW, halfH float64) {
	maxSize := math.Max(common.SwordAttachmentWidth, common.SwordAttachmentHeight) / 2
	return b.Position.Px, b.Position.Py, maxSize, maxSize
}

// Для OBB (точный)
func (b *BlueSword) GetOBB() (centerX, centerY, halfW, halfH, angle float64) {
	return b.Position.Px, b.Position.Py,
		common.SwordAttachmentWidth / 2,
		common.SwordAttachmentHeight / 2,
		b.Angle * math.Pi / 180
}

// GetHitBoxID возвращает уникальный ID для идентификации
func (b *BlueSword) GetHitBoxID() string {
	return b.Tag()
}

// IsStatic проверяет, статичен ли объект (стена, платформа)
// Если true - объект не двигается при отталкивании
func (b *BlueSword) IsStatic() bool {
	return false
}

// ApplyPush применяет силу отталкивания (сдвиг)
func (b *BlueSword) ApplyPush(x, y float64) {

}

// OnCollision вызывается при столкновении с другим объектом
func (b *BlueSword) OnCollision(other hitboxes.HitBoxer) {}

func (b *BlueSword) Tag() string {
	return "blueSword"
}

func (b *BlueSword) IsActive() bool {
	return true
}
