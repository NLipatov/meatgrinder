package domain

import "fmt"

type Warrior struct {
	BaseCharacter

	attackPower  float64
	attackRadius float64
	physicalRes  float64
}

func NewWarrior(id string, x, y float64) *Warrior {
	w := &Warrior{
		attackPower:  20,
		attackRadius: 5.0,
		physicalRes:  0.5,
	}
	w.InitBase(id, 100, x, y, 1.0, Physical)
	return w
}

func (w *Warrior) Attack(targets []Character) {
	for _, t := range targets {
		if t.IsDead() {
			continue
		}
		damage := w.attackPower
		if otherW, ok := t.(*Warrior); ok {
			damage *= 1.0 - otherW.physicalRes
		}
		t.TakeDamage(damage, Physical)
		fmt.Printf("Warrior %s dealt %f damage to %s\n", w.ID(), damage, t.ID())
	}
}

func (w *Warrior) TakeDamage(amount float64, dmgType DamageType) {
	if dmgType == Physical {
		amount *= 1.0 - w.physicalRes
	}
	w.BaseCharacter.TakeDamage(amount, dmgType)
}

func (w *Warrior) AttackPower() float64  { return w.attackPower }
func (w *Warrior) AttackRadius() float64 { return w.attackRadius }
