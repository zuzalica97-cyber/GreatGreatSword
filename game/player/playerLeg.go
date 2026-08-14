package player

import (
	"great-sword/game"
	playerabilities "great-sword/game/abilities/playerAbilities"
	"great-sword/game/common"
	"great-sword/game/hitboxes"
	gameL "great-sword/game/world"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

var _ game.Entity = (*PlayerLeg)(nil)
var _ game.PlayerLegInter = (*PlayerLeg)(nil)
var _ hitboxes.HitBoxer = (*PlayerLeg)(nil)

type PlayerLeg struct {
	Position         common.PointPlayer
	Speed            common.PointSpeed
	Texture          *ebiten.Image
	AxelerationLeg   float64
	DeAxelerationLeg float64
	CurrentMaxSpeed  float64
	Rebound          float64
	MoveX            float64
	MoveY            float64
	Weight           float64
	Density          float64

	AbilityLegManager *gameL.PlayerWorld
}

func NewPlayerLeg(manager *hitboxes.CollisionManager) *PlayerLeg {

	p := &PlayerLeg{
		Position: common.PointPlayer{
			Px: common.RoomWidth/2 - common.PlayerSize/2,
			Py: common.RoomHeight/2 - common.PlayerSize/2,
		},
		Speed: common.PointSpeed{
			Vx: 0,
			Vy: 0,
		},
		Weight:            4,
		Density:           2.0,
		AbilityLegManager: gameL.NewPlayerWorld(),
	}
	p.AbilityLegManager.AddAbility(
		playerabilities.NewDash(),
	)

	if manager != nil {
		manager.AddObject(p)
	}

	return p
}

func (p *PlayerLeg) ActivateBoost() {
	BoostTimer = BoostTimerLong
}

func (p *PlayerLeg) Update(worldView game.WorldView, manager *hitboxes.CollisionManager) bool {

	dt := 1.0 / 60.0

	common.SwordExist = SwordIxist

	if common.PlayerHelth >= common.MaxPlayerHelth {
		common.PlayerHelth = common.MaxPlayerHelth
	}

	if BoostTimer > 0 {
		BoostTimer -= dt
	}

	if SwordIxistTimer > 0 {
		SwordIxistTimer -= dt
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.AbilityLegManager.ActivateAbility("Dash", worldView)
	}

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		SwordVanished()
	}

	p.Rebound = Rebount

	p.AxelerationLeg = common.Acceleration
	p.DeAxelerationLeg = common.Deceleration

	p.CurrentMaxSpeed = common.MaxSpeed

	if BoostTimer > 0 {
		p.CurrentMaxSpeed = common.MaxSpeed + float64(BoostSpeed)
	}
	if !SwordIxist {
		p.CurrentMaxSpeed = common.MaxSpeed * 2
		p.DeAxelerationLeg = common.Deceleration / 5
		p.AxelerationLeg = common.Acceleration * 2

	}

	if p.CurrentMaxSpeed > float64(common.MaxPlayerSpeedMoving) {
		p.CurrentMaxSpeed = float64(common.MaxPlayerSpeedMoving)
	}

	moveX, moveY := 0.0, 0.0

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		moveX = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		moveX = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		moveY = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		moveY = 1
	}

	if moveX != 0 && moveY != 0 { //движение подиоганале
		moveX *= 0.7071
		moveY *= 0.7071
	}

	if !SwordIxist {
		p.CurrentMaxSpeed = common.MaxSpeed * 1.5
		p.AxelerationLeg = common.Acceleration * 1.8   // резкий разгон
		p.DeAxelerationLeg = common.Deceleration * 2.0 // резкое торможение
	}

	p.MoveX = moveX
	p.MoveY = moveY

	p.AbilityLegManager.UpdateAbilities(worldView)

	if p.MoveX != 0 { //Применяем ускорение
		p.Speed.Vx += p.MoveX * p.AxelerationLeg * dt
	} else {
		dec := p.DeAxelerationLeg * dt // Замедление если ненажата клавиша
		if math.Abs(p.Speed.Vx) > dec {
			p.Speed.Vx -= math.Copysign(dec, p.Speed.Vx)
		} else {
			p.Speed.Vx = 0
		}
	}

	if p.MoveY != 0 { //Применяем ускорение
		p.Speed.Vy += p.MoveY * p.DeAxelerationLeg * dt
	} else {
		dec := p.DeAxelerationLeg * dt // Замедление если ненажата клавиша
		if math.Abs(p.Speed.Vy) > dec {
			p.Speed.Vy -= math.Copysign(dec, p.Speed.Vy)
		} else {
			p.Speed.Vy = 0
		}
	}

	if p.Speed.Vx > p.CurrentMaxSpeed {
		p.Speed.Vx = p.CurrentMaxSpeed
	}
	if p.Speed.Vx < -p.CurrentMaxSpeed {
		p.Speed.Vx = -p.CurrentMaxSpeed
	}
	if p.Speed.Vy > p.CurrentMaxSpeed {
		p.Speed.Vy = p.CurrentMaxSpeed
	}
	if p.Speed.Vy < -p.CurrentMaxSpeed {
		p.Speed.Vy = -p.CurrentMaxSpeed
	}

	p.Position.Px += p.Speed.Vx * dt
	p.Position.Py += p.Speed.Vy * dt

	//Границы с отскоком и потерей скорости
	if p.Position.Px < 0 {
		p.Position.Px = 0
		p.Speed.Vx = -p.Speed.Vx * p.Rebound //отскок с потерей скорости
	}
	if p.Position.Px > common.RoomWidth-common.PlayerSize {
		p.Position.Px = common.RoomHeight - common.PlayerSize
		p.Speed.Vx = -p.Speed.Vx * p.Rebound //отскок с потерей скорости
	}
	if p.Position.Py < 0 {
		p.Position.Py = 0
		p.Speed.Vy = -p.Speed.Vy * p.Rebound //отскок с потерей скорости
	}
	if p.Position.Py > common.RoomHeight-common.PlayerSize {
		p.Position.Py = common.RoomHeight - common.PlayerSize
		p.Speed.Vy = -p.Speed.Vy * p.Rebound //отскок с потерей скорости
	}

	return false
}

func (p *PlayerLeg) GetAABB() (posX, posY, halfW, halfH float64) {
	halfSize := common.PlayerSize / 2
	return p.Position.Px + float64(halfSize), p.Position.Py + float64(halfSize), float64(halfSize), float64(halfSize)
}

// GetHitBoxID возвращает уникальный ID для идентификации
func (p *PlayerLeg) GetHitBoxID() string {
	return p.Tag()
}

// IsStatic проверяет, статичен ли объект (стена, платформа)
// Если true - объект не двигается при отталкивании
func (p *PlayerLeg) IsStatic() bool {
	return false
}

// ApplyPush применяет силу отталкивания (сдвиг)
func (p *PlayerLeg) ApplyPush(x, y float64) {
	p.Position.Px += x
	p.Position.Py += y

}

// GetWeight - возвращает вес объекта
func (b *PlayerLeg) GetWeight() float64 {
	return b.Weight
}

// GetDensity - возвращает плотность объекта
func (b *PlayerLeg) GetDensity() float64 {
	return b.Density
}

func (p *PlayerLeg) HasAura() bool {
	return false
}

func (p *PlayerLeg) AffectedByAura() bool {
	return false
}

// OnCollision вызывается при столкновении с другим объектом
func (p *PlayerLeg) OnCollision(other hitboxes.HitBoxer) {}

func (p *PlayerLeg) Draw(screen *ebiten.Image, camera *kamera.Camera) {
}

func (p *PlayerLeg) Tag() string {
	return "playerLeg"
}

func (p *PlayerLeg) IsActive() bool {
	return true
}
