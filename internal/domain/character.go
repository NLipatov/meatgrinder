package domain

import "math"

type Character interface {
	ID() string
	Position() (float64, float64)
	Health() float64
	IsDead() bool
	DamageType() DamageType
	Attack(targets []Character)
	MoveTo(x, y float64)
	TakeDamage(amount float64, damageType DamageType)
	AttackPower() float64
	AttackRadius() float64

	Update(dt float64)
}

type DamageType int

const (
	Unset DamageType = iota
	Physical
	Magical
)

type BaseCharacter struct {
	id         string
	health     float64
	x, y       float64
	baseSpeed  float64
	speed      float64
	isDead     bool
	damageType DamageType

	slowTimer  float64
	slowAmount float64
}

func (bc *BaseCharacter) InitBase(id string, health float64, x, y float64, baseSpeed float64, damageType DamageType) {
	bc.id = id
	bc.health = health
	bc.x = x
	bc.y = y
	bc.baseSpeed = baseSpeed
	bc.speed = baseSpeed
	bc.damageType = damageType
}

func (bc *BaseCharacter) Update(dt float64) {
	if bc.slowTimer > 0 {
		bc.slowTimer -= dt
		if bc.slowTimer <= 0 {
			bc.speed = bc.baseSpeed
			bc.slowTimer = 0
			bc.slowAmount = 0
		}
	}
}
func (bc *BaseCharacter) ID() string                   { return bc.id }
func (bc *BaseCharacter) Position() (float64, float64) { return bc.x, bc.y }
func (bc *BaseCharacter) Health() float64              { return bc.health }
func (bc *BaseCharacter) IsDead() bool                 { return bc.isDead }

func (bc *BaseCharacter) MoveTo(x, y float64) {
	dx := x - bc.x
	dy := y - bc.y
	dist := math.Hypot(dx, dy)
	if dist > 0 {
		bc.x += (dx / dist) * bc.speed
		bc.y += (dy / dist) * bc.speed
	}
}

func (bc *BaseCharacter) TakeDamage(amount float64, dmgType DamageType) {
	bc.health -= amount
	if bc.health <= 0 {
		bc.isDead = true
	}
}

func (bc *BaseCharacter) Attack(targets []Character) {}

func (bc *BaseCharacter) ApplySlow(amount float64, duration float64) {
	if bc.isDead {
		return
	}

	bc.slowAmount = amount
	bc.slowTimer = duration
	bc.speed = bc.baseSpeed * (1.0 - bc.slowAmount)
}

func (bc *BaseCharacter) DamageType() DamageType {
	return bc.damageType
}
