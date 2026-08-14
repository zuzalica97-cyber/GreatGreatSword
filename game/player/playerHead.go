package player

import (
	"great-sword/game"
	playerabilities "great-sword/game/abilities/playerAbilities"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
	"image/color"
	"log"

	"math"

	gameH "great-sword/game/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/setanarut/kamera/v2"
)

var _ game.Entity = (*PlayerHead)(nil)
var _ game.PlayerHeadInter = (*PlayerHead)(nil)

var SpinRotatinInput int

type PlayerHead struct {
	PositionHead    common.PointPlayer
	Angle           float64
	AngularVelocity float64

	BoostTimer         float64
	BoostActive        bool
	NormalSpeed        float64
	NormalROTSPEED     float64
	CurrentRotSpead    float64
	HDeceleration      float64
	HAceleration       float64
	Texture            *ebiten.Image
	AbilityHeadManager *gameH.PlayerWorld
}

func NewPlayerHead() *PlayerHead {
	p := &PlayerHead{
		Angle:              0.0,
		AngularVelocity:    0.0,
		AbilityHeadManager: gameH.NewPlayerWorld(),
	}

	var err error

	p.Texture, _, err = ebitenutil.NewImageFromFile("assets/playerV0.2.png") //НЕ РАБОТАЕТ СДЕСЬ
	if err != nil {
		log.Fatal("failed to load lower texture", err)
	}

	p.AbilityHeadManager.AddAbility(
		playerabilities.NewFastRotate(),
	)

	return p
}

func (p *PlayerHead) Update(woroldWiev game.WorldView, manager *hitboxes.CollisionManager) bool {
	dt := 1.0 / 60.0

	p.CurrentRotSpead = common.HeadRotationSpeed
	p.HDeceleration = common.HeadDeceleration
	p.HAceleration = common.HeadAcceleration

	if BoostTimer > 0 {
		p.CurrentRotSpead = common.HeadRotationSpeed + float64(BoostRotating)
	}

	rotationInput := 0.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		rotationInput = -1
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		rotationInput = 1
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) &&
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		rotationInput = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyE) {
		p.AbilityHeadManager.ActivateAbility("FastRotate", woroldWiev)
	}

	p.AbilityHeadManager.UpdateAbilities(woroldWiev)

	if rotationInput != 0 {
		p.AngularVelocity += rotationInput * p.HAceleration * dt
	} else {
		dec := p.HDeceleration * dt
		if math.Abs(p.AngularVelocity) > dec {
			p.AngularVelocity -= math.Copysign(dec, p.AngularVelocity)
		} else {
			p.AngularVelocity = 0
		}
	}

	if p.AngularVelocity > p.CurrentRotSpead {
		p.AngularVelocity = p.CurrentRotSpead
	}
	if p.AngularVelocity < -p.CurrentRotSpead {
		p.AngularVelocity = -p.CurrentRotSpead
	}

	p.Angle += p.AngularVelocity * dt

	p.Angle = math.Mod(p.Angle, 360)
	if p.Angle < 0 {
		p.Angle += 360
	}

	for _, entity := range woroldWiev.SearchEntities("playerLeg") {
		pL := entity.(*PlayerLeg)

		centerX := pL.Position.Px + common.PlayerSize/2
		centerY := pL.Position.Py + common.PlayerSize/2

		p.PositionHead.Px = centerX
		p.PositionHead.Py = centerY

		vx := pL.Speed.Vx
		vy := pL.Speed.Vy

		if vx < 0 {
			vx = -vx
		}
		if vy < 0 {
			vy = -vy
		}
		if p.CurrentRotSpead < 0 {
			p.CurrentRotSpead = -p.CurrentRotSpead
		}

		FinalDamage := ((pL.Speed.Vx + pL.Speed.Vy) * 0.1) + (p.CurrentRotSpead * 0.1) + 50

		if p.BoostActive {
			FinalDamage = FinalDamage * 1.5
		}

		common.PlayerDamage = FinalDamage
	}

	return false
}

func (p *PlayerHead) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	// Центр игрока в мировых координатах
	centerX := p.PositionHead.Px + common.PlayerHeadSize/2
	centerY := p.PositionHead.Py + common.PlayerHeadSize/2

	// Экранные координаты с учётом камеры
	screenX := centerX - camera.X
	screenY := centerY - camera.Y

	// 1. Рисуем повёрнутый красный квадрат
	DrawRotatedRect(
		screen,
		screenX,
		screenY,
		common.PlayerHeadSize,
		common.PlayerHeadSize,
		p.Angle,                    // ← угол поворота
		color.RGBA{255, 0, 0, 255}, // ярко-красный
	)
}

func (p *PlayerHead) Tag() string {
	return "playerHead"
}

func (p *PlayerHead) IsActive() bool {
	return true
}
