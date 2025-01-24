package domain

import "fmt"

type Warrior struct {
	BaseCharacter
	power  float64
	radius float64
	res    float64
}

func NewWarrior(id string, x, y float64) *Warrior {
	return &Warrior{
		BaseCharacter: BaseCharacter{
			id:         id,
			health:     100,
			x:          x,
			y:          y,
			speed:      1,
			damageType: Physical,
			state:      StateIdle,
		},
		power:  20,
		radius: 250,
		res:    0.5,
	}
}

func (w *Warrior) Attack(targets []Character) {
	if w.isDead || w.state == StateDying {
		return
	}
	w.state = StateAttacking
	w.attackTimer = 0.3
	for _, t := range targets {
		if t.IsDead() {
			continue
		}
		d := w.power
		if ww, ok := t.(*Warrior); ok {
			d *= (1 - ww.res)
		}
		t.TakeDamage(d, Physical)
		fmt.Printf("%s warrior attacked %s\n", w.id, t.ID())
	}
}

func (w *Warrior) TakeDamage(a float64, dt DamageType) {
	if dt == Physical {
		a *= (1 - w.res)
	}
	w.BaseCharacter.TakeDamage(a, dt)
}

func (w *Warrior) AttackPower() float64  { return w.power }
func (w *Warrior) AttackRadius() float64 { return w.radius }
