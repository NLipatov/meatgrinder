package domain_test

import (
	"meatgrinder/internal/domain"
	"testing"
)

func TestWarrior_Attack(t *testing.T) {
	w := domain.NewWarrior("w1", 0, 0)
	m := domain.NewMage("m1", 0, 0)
	targets := []domain.Character{m}

	w.Attack(targets)
	if m.Health() >= 80 {
		t.Errorf("Expected mage health < 80 after attack, got %v", m.Health())
	}
}
