package domain

import "math"

type Character interface {
	ID() string
	Position() (float64, float64)
	Health() float64
	IsDead() bool

	Attack(targets []Character)
	MoveTo(x, y float64)
	TakeDamage(amount float64, damageType DamageType)

	Update(dt float64)
}

type DamageType int

const (
	Physical DamageType = iota
	Magical
)

type BaseCharacter struct {
	id     string
	health float64
	x, y   float64
	speed  float64
	isDead bool
}

func NewBaseCharacter(health, x, y, speed float64, isDead bool) BaseCharacter {
	return BaseCharacter{
		id:     "",
		health: health,
		x:      x,
		y:      y,
		speed:  speed,
		isDead: isDead,
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
func (bc *BaseCharacter) Update(dt float64)          {}
