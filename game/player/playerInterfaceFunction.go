package player

import (
	"math"
)

func (p *PlayerLeg) ApplyForce(vx, vy float64) {
	p.Speed.Vx = vx
	p.Speed.Vy = vy
}

func (p *PlayerLeg) GetSpeedXY() (float64, float64) {
	return p.Speed.Vx, p.Speed.Vy
}

func (p *PlayerLeg) SetSpeedXY(vx, vy float64) {
	p.Speed.Vx = vx
	p.Speed.Vy = vy
}

func (p *PlayerLeg) GetMaxSpeed() float64 {
	return p.CurrentMaxSpeed
}

func (p *PlayerLeg) SetMaxSpeed(speed float64) {
	p.CurrentMaxSpeed = speed
}

func (p *PlayerLeg) GetDirection() (float64, float64) {
	return p.MoveX, p.MoveY
}

func (p *PlayerLeg) SetDirection(x, y float64) {
	p.MoveX = x
	p.MoveY = y
}

func (p *PlayerLeg) GetAxelDexel() (float64, float64) {
	return p.AxelerationLeg, p.DeAxelerationLeg
}

func (p *PlayerLeg) SetAxelDexel(Axel, Dexel float64) {
	p.AxelerationLeg = Axel
	p.DeAxelerationLeg = Dexel
}

// ---- ПОЗИЦИЯ ----
func (p *PlayerLeg) GetPosition() (float64, float64) {
	return p.Position.Px, p.Position.Py
}

func (p *PlayerLeg) SetPosition(x, y float64) {
	p.Position.Px = x
	p.Position.Py = y
}

// ---- СОСТОЯНИЕ ----
func (p *PlayerLeg) IsMoving() bool {
	return p.Speed.Vx != 0 || p.Speed.Vy != 0
}

func (p *PlayerLeg) IsGrounded() bool {
	return p.Speed.Vy == 0
}

// ---- ВСПОМОГАТЕЛЬНЫЕ ----
func (p *PlayerLeg) GetLeg() interface{} {
	return p
}

func (p *PlayerLeg) GetHead() interface{} {
	// Ищем Head через world
	// for _, entity := range p.world.SearchEntities("playerHead") {
	//     return entity
	// }
	return nil
}

//ГОЛОВА

// ============================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА game.Player ДЛЯ PlayerHead
// ============================================================

// ---- ВРАЩЕНИЕ ----
func (p *PlayerHead) GetRotationSpeed() float64 {
	return p.AngularVelocity
}

func (p *PlayerHead) SetRotationSpeed(speed float64) {
	p.AngularVelocity = speed
}

func (p *PlayerHead) GetAngle() float64 {
	return p.Angle
}

func (p *PlayerHead) SetAngle(angle float64) {
	p.Angle = angle
}

// ---- ПОЗИЦИЯ ----
func (p *PlayerHead) GetPosition() (float64, float64) {
	return p.PositionHead.Px, p.PositionHead.Py
}

func (p *PlayerHead) SetPosition(x, y float64) {
	p.PositionHead.Px = x
	p.PositionHead.Py = y
}

// ---- ДВИЖЕНИЕ (по углу головы) ----
func (p *PlayerHead) GetDirection() (float64, float64) {
	angleRad := p.Angle * math.Pi / 180
	return math.Cos(angleRad), math.Sin(angleRad)
}

func (p *PlayerHead) SetCurrentRotSpeed(speed float64) {
	p.CurrentRotSpead = speed
}

func (p *PlayerHead) GetAceleration() float64 {
	return p.HAceleration
}
func (p *PlayerHead) SetAceleration(axel float64) {
	p.HAceleration = axel
}
func (p *PlayerHead) GetDeAceleration() float64 {
	return p.HDeceleration
}
func (p *PlayerHead) SetDeAceleration(dexel float64) {
	p.HDeceleration = dexel
}

// ---- СОСТОЯНИЕ ----
func (p *PlayerHead) IsMoving() bool {
	return p.AngularVelocity != 0
}

func (p *PlayerHead) IsGrounded() bool {
	return true
}

// ---- ВСПОМОГАТЕЛЬНЫЕ ----

func (p *PlayerHead) GetHead() interface{} {
	return p
}

func (p *PlayerHead) GetCurrentRotSpeed() float64 {
	return p.CurrentRotSpead
}
