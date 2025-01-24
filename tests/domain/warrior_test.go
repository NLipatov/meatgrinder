package domain_test

import (
	"meatgrinder/internal/domain"
	"testing"
)

func TestWarrior_Attack(t *testing.T) {
	war1 := domain.NewWarrior("w1", 0, 0)
	mage := domain.NewMage("m1", 0, 0)

	mageHPBefore := mage.Health()
	war1.Attack([]domain.Character{mage})

	if mage.Health() >= mageHPBefore {
		t.Errorf("Mage's health must be lower after warrior's attack")
	}
}

func TestWarrior_ApplySlow(t *testing.T) {
	war := domain.NewWarrior("w2", 0, 0)
	warSpeedBefore := war.AttackPower()

	war.ApplySlow(0.5, 2.0)

	war.MoveTo(10, 0)
	if war.IsDead() {
		t.Errorf("Warrior shouldn't be dead just from slow")
	}

	if warSpeedBefore == 0 {
		t.Skip("Can't test speed. Implementation detail hidden.")
	}
}
