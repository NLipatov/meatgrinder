package domain

import "fmt"

type Mage struct {
	BaseCharacter
	power  float64
	radius float64
	res    float64
}

func NewMage(id string, x, y float64) *Mage {
	return &Mage{
		BaseCharacter: BaseCharacter{
			id:         id,
			health:     80,
			x:          x,
			y:          y,
			speed:      7,
			damageType: Magical,
			state:      StateIdle,
		},
		power:  30,
		radius: 450,
		res:    0.3,
	}
}

func (m *Mage) Attack(targets []Character) {
	if m.isDead || m.state == StateDying {
		return
	}
	m.state = StateAttacking
	m.attackTimer = 0.3
	for _, t := range targets {
		if t.IsDead() {
			continue
		}
		t.TakeDamage(m.power, Magical)
		fmt.Printf("%s mage attacked %s\n", m.id, t.ID())
	}
}

func (m *Mage) TakeDamage(a float64, dt DamageType) {
	if dt == Magical {
		a *= (1 - m.res)
	}
	m.BaseCharacter.TakeDamage(a, dt)
}

func (m *Mage) AttackPower() float64  { return m.power }
func (m *Mage) AttackRadius() float64 { return m.radius }
