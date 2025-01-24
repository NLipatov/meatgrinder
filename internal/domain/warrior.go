package domain

import "fmt"

type Warrior struct {
	BaseCharacter
	attackPower  float64
	attackRadius float64
	physicalRes  float64
}

func NewWarrior(id string, x, y float64) *Warrior {
	return &Warrior{
		BaseCharacter: BaseCharacter{
			id:         id,
			health:     100,
			x:          x,
			y:          y,
			speed:      1.0,
			damageType: Physical,
		},
		attackPower:  20,
		attackRadius: 5.0,
		physicalRes:  0.5,
	}
}

func (w *Warrior) Attack(targets []Character) {
	for _, t := range targets {
		if t.IsDead() {
			continue
		}
		tx, ty := t.Position()
		wx, wy := w.Position()
		dist := distance(wx, wy, tx, ty)
		if dist <= w.attackRadius {
			var dmg float64 = w.attackPower
			if t, ok := t.(*Warrior); ok {
				dmg *= 1.0 - t.physicalRes
			}
			t.TakeDamage(dmg, Physical)
			fmt.Printf("Warrior %s dealt %f damage to %s\n", w.ID(), dmg, t.ID())
		}
	}
}

func (w *Warrior) TakeDamage(amount float64, dmgType DamageType) {
	if dmgType == Physical {
		amount *= 1.0 - w.physicalRes
	}
	w.BaseCharacter.TakeDamage(amount, dmgType)
}

func (w *Warrior) AttackRadius() float64 {
	return w.attackRadius
}

func (w *Warrior) AttackPower() float64 {
	return w.attackPower
}

func distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}
