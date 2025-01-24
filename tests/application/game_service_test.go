package application

import (
	"context"
	"fmt"
	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/application/services"
	"meatgrinder/internal/domain"
	"strings"
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

func TestGameService_ProcessMoveCommand(t *testing.T) {
	w := domain.NewWorld(100, 100)
	l := &mockLogger{}
	gs := services.NewGameService(w, l)

	w.SpawnRandomCharacter("char1")

	dto := dtos.CommandDTO{
		Type:        "MOVE",
		CharacterID: "char1",
		Data: map[string]interface{}{
			"x": 5.0,
			"y": 10.0,
		},
	}

	err := gs.ProcessCommandDTO(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ch := w.Characters["char1"]
	x, y := ch.Position()
	if x == 0 && y == 0 {
		t.Error("character did not move (still at 0,0)")
	}

	if len(l.logs) == 0 {
		t.Error("no log message recorded")
	}
}

func TestGameService_MageAttackWarrior(t *testing.T) {
	w := domain.NewWorld(100, 100)
	l := &mockLogger{}
	gs := services.NewGameService(w, l)

	mage := domain.NewMage("Mage", 20, 30)
	warrior := domain.NewWarrior("Warrior", 20, 30)
	w.Characters[mage.ID()] = mage
	w.Characters[warrior.ID()] = warrior
	mageInitHealth := w.Characters["Mage"].Health()
	warriorInitHealth := w.Characters["Warrior"].Health()

	attackDTO := dtos.CommandDTO{
		Type:        "ATTACK",
		CharacterID: "Mage",
		Data: map[string]interface{}{
			"target_id": "Warrior",
		},
	}

	err := gs.ProcessCommandDTO(attackDTO)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, msg := range l.logs {
		if msg == "Mage Mage attacked Warrior with magic" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected attack log not found, messages: %v", l.logs)
	}

	if w.Characters["Mage"].Health() != mageInitHealth {
		t.Fatalf("Expected Mage health to be %v, got %v", mageInitHealth, w.Characters["Mage"].Health())
	}

	if w.Characters["Warrior"].Health() == warriorInitHealth {
		t.Fatalf("Expected Warrior health to be less than %v, got %v", warriorInitHealth, w.Characters["Warrior"].Health())
	}
}

func TestGameService_WarriorAttackMage(t *testing.T) {
	w := domain.NewWorld(100, 100)
	l := &mockLogger{}
	gs := services.NewGameService(w, l)

	mage := domain.NewMage("Mage", 20, 30)
	warrior := domain.NewWarrior("Warrior", 20, 30)
	w.Characters[mage.ID()] = mage
	w.Characters[warrior.ID()] = warrior
	mageInitHealth := w.Characters["Mage"].Health()
	warriorInitHealth := w.Characters["Warrior"].Health()

	attackDTO := dtos.CommandDTO{
		Type:        "ATTACK",
		CharacterID: "Warrior",
		Data: map[string]interface{}{
			"target_id": "Mage",
		},
	}

	err := gs.ProcessCommandDTO(attackDTO)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, msg := range l.logs {
		if msg == fmt.Sprintf("Warrior %s did an AoE melee attack", "Warrior") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected attack log not found, messages: %v", l.logs)
	}

	if w.Characters["Mage"].Health() == mageInitHealth {
		t.Fatalf("Expected Mage health to be less than %v, got %v", mageInitHealth, w.Characters["Mage"].Health())
	}

	if w.Characters["Warrior"].Health() != warriorInitHealth {
		t.Fatalf("Expected Warrior health to be equal to %v, got %v", warriorInitHealth, w.Characters["Warrior"].Health())
	}
}

func TestGameService_WarriorAoEAttack(t *testing.T) {
	w := domain.NewWorld(100, 100)
	l := &mockLogger{}
	gs := services.NewGameService(w, l)

	warrior := domain.NewWarrior("war1", 0, 0)
	m1 := domain.NewMage("mageA", 1, 0)
	m2 := domain.NewMage("mageB", 3, 0)
	w.Characters["war1"] = warrior
	w.Characters["mageA"] = m1
	w.Characters["mageB"] = m2

	attackDTO := dtos.CommandDTO{
		Type:        "ATTACK",
		CharacterID: "war1",
		Data:        map[string]interface{}{},
	}

	err := gs.ProcessCommandDTO(attackDTO)
	if err != nil {
		t.Fatalf("attack error: %v", err)
	}

	if m1.Health() == 80 {
		t.Errorf("mageA should have taken damage, got health=%.1f", m1.Health())
	}
	if m2.Health() == 80 {
		t.Errorf("mageB should have taken damage, got health=%.1f", m2.Health())
	}

	found := false
	for _, msg := range l.logs {
		if strings.Contains(msg, "Warrior war1 did an AoE melee attack") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected AoE log not found")
	}
}
