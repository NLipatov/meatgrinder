package domain_test

import (
	"meatgrinder/internal/domain"
	"testing"
)

func TestMage_Attack(t *testing.T) {
	mage := domain.NewMage("mage1", 10, 10)
	warr := domain.NewWarrior("war1", 10, 10)

	warrHealthBefore := warr.Health()

	mage.Attack([]domain.Character{warr})

	if warr.Health() >= warrHealthBefore {
		t.Errorf("Warrior's health should decrease after Mage attack. Before=%.1f, After=%.1f",
			warrHealthBefore, warr.Health())
	}
}

func TestMage_TakeDamage_MagicalRes(t *testing.T) {
	mage := domain.NewMage("mage1", 0, 0)
	hpBefore := mage.Health()

	mage.TakeDamage(50, domain.Magical)
	expectedHP := hpBefore - (50 * (1.0 - 0.3))

	if mage.Health() != expectedHP {
		t.Errorf("Expected Mage HP=%.1f after magical damage, got %.1f", expectedHP, mage.Health())
	}
}
