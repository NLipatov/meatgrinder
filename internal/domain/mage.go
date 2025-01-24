package domain

import "fmt"

type Mage struct {
	BaseCharacter

	attackPower  float64
	magicalRes   float64
	attackRadius float64
}

func NewMage(id string, x, y float64) *Mage {
	m := &Mage{
		attackPower:  30,
		magicalRes:   0.3,
		attackRadius: 8.0,
	}
	m.InitBase(id, 80, x, y, 1.5, Magical)
	return m
}

func (m *Mage) Attack(targets []Character) {
	for _, t := range targets {
		if t.IsDead() {
			continue
		}
		t.TakeDamage(m.attackPower, Magical)
		fmt.Printf("Mage %s dealt %f magical damage to %s\n", m.ID(), m.attackPower, t.ID())
	}
}

func (m *Mage) TakeDamage(amount float64, dmgType DamageType) {
	if dmgType == Magical {
		amount *= 1.0 - m.magicalRes
	}
	m.BaseCharacter.TakeDamage(amount, dmgType)
}

func (m *Mage) AttackPower() float64  { return m.attackPower }
func (m *Mage) AttackRadius() float64 { return m.attackRadius }
