package domain

import "math"

type DamageType int

const (
	Unset DamageType = iota
	Physical
	Magical
)

type CharacterState string

const (
	StateIdle      CharacterState = "idle"
	StateRunning   CharacterState = "running"
	StateAttacking CharacterState = "attacking"
	StateDying     CharacterState = "dying"
)

type Character interface {
	ID() string
	Position() (float64, float64)
	Health() float64
	IsDead() bool
	DamageType() DamageType
	State() CharacterState
	SetState(CharacterState)
	Attack([]Character)
	MoveStep(float64, float64)
	TakeDamage(float64, DamageType)
	AttackPower() float64
	AttackRadius() float64
	Update(float64)
	FlashRed() bool
}

type BaseCharacter struct {
	id          string
	health      float64
	x, y        float64
	isDead      bool
	damageType  DamageType
	speed       float64
	state       CharacterState
	attackTimer float64
	hitTimer    float64
	flashRedOn  bool
	noMoveTimer float64
}

func (bc *BaseCharacter) ID() string                   { return bc.id }
func (bc *BaseCharacter) Position() (float64, float64) { return bc.x, bc.y }
func (bc *BaseCharacter) Health() float64              { return bc.health }
func (bc *BaseCharacter) IsDead() bool                 { return bc.isDead }
func (bc *BaseCharacter) DamageType() DamageType       { return bc.damageType }
func (bc *BaseCharacter) State() CharacterState        { return bc.state }
func (bc *BaseCharacter) SetState(s CharacterState)    { bc.state = s }
func (bc *BaseCharacter) FlashRed() bool               { return bc.flashRedOn }

func (bc *BaseCharacter) MoveStep(dx, dy float64) {
	if bc.isDead || bc.state == StateDying {
		return
	}
	dist := math.Hypot(dx, dy)
	if dist < 0.0001 {
		return
	}
	bc.x += (dx / dist) * bc.speed
	bc.y += (dy / dist) * bc.speed
	bc.state = StateRunning
	bc.noMoveTimer = 0
}

func (bc *BaseCharacter) TakeDamage(amt float64, dt DamageType) {
	if bc.isDead || bc.state == StateDying {
		return
	}
	bc.health -= amt
	bc.flashRedOn = true
	bc.hitTimer = 0.2
	if bc.health <= 0 {
		bc.isDead = true
		bc.state = StateDying
	}
}

func (bc *BaseCharacter) Attack([]Character)    {}
func (bc *BaseCharacter) AttackPower() float64  { return 0 }
func (bc *BaseCharacter) AttackRadius() float64 { return 0 }

func (bc *BaseCharacter) Update(dt float64) {
	if bc.isDead && bc.state != StateDying {
		bc.state = StateDying
	}

	if bc.state == StateAttacking {
		bc.attackTimer -= dt
		if bc.attackTimer <= 0 {
			bc.state = StateIdle
			bc.attackTimer = 0
		}
	}

	if bc.hitTimer > 0 {
		bc.hitTimer -= dt
		if bc.hitTimer <= 0 {
			bc.flashRedOn = false
			bc.hitTimer = 0
		}
	}

	if bc.state == StateRunning {
		bc.noMoveTimer += dt
		if bc.noMoveTimer >= 0.5 {
			bc.state = StateIdle
		}
	}
}
