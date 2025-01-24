package domain_test

import (
	"meatgrinder/internal/domain"
	"testing"
)

func TestBaseCharacter_TakeDamage(t *testing.T) {
	health := 100.0
	ch := domain.NewBaseCharacter(health, 1.1, 2.2, 1, false)

	ch.TakeDamage(50, domain.Physical)
	if ch.Health() != 50.0 {
		t.Errorf("Expected health 50, got %v", ch.Health())
	}
}
