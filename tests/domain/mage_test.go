package domain_test

import (
	"meatgrinder/internal/domain"
	"testing"
)

func TestMage_Attack(t *testing.T) {
	w := domain.NewWarrior("w1", 0, 0)
	m := domain.NewMage("m1", 0, 0)
	m.Attack([]domain.Character{w})

	if w.Health() >= 100 {
		t.Errorf("Warrior health must be < 100 after mage's attack")
	}
}
