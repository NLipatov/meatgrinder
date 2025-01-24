package application

import (
	"context"
	"meatgrinder/internal/application/services"
	"meatgrinder/internal/domain"
	"testing"
)

func TestGameService_BroadcastState_Attack(t *testing.T) {
	w := domain.NewWorld(100, 100)
	gs := services.NewGameService(w, &mockLogger{})
	ctx := context.Background()
	defer ctx.Done()

	go gs.BroadcastState(ctx)

	mage := domain.NewMage("Mage", 20, 30)
	warrior := domain.NewMage("Warrior", 20, 30)
	w.Characters[mage.ID()] = mage
	w.Characters[warrior.ID()] = warrior

	warriorTargets := make([]domain.Character, 1)
	warriorTargets[0] = mage

	mageTargets := make([]domain.Character, 1)
	mageTargets[0] = warrior

	w.Characters[warrior.ID()].Attack(warriorTargets)
	w.Characters[mage.ID()].Attack(mageTargets)

	if w.Characters[mage.ID()].Health() == domain.NewMage("", 1, 1).Health() {
		t.Errorf("Expected mage health to be less than %v, got %v",
			domain.NewMage("", 1, 1).Health(), w.Characters[mage.ID()].Health())
	}

	if w.Characters[warrior.ID()].Health() == domain.NewWarrior("", 1, 1).Health() {
		t.Errorf("Expected warrior health to be less than %v, got %v",
			domain.NewWarrior("", 1, 1).Health(), w.Characters[warrior.ID()].Health())
	}
}
