package application_test

import (
	"context"
	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/application/services"
	"meatgrinder/internal/domain"
	"testing"
)

type mockLogger struct {
	events []string
}

func (ml *mockLogger) LogEvent(e string) {
	ml.events = append(ml.events, e)
}

func TestGameService_SpawnMoveAttack(t *testing.T) {
	w := domain.NewWorld(100, 100)
	logger := &mockLogger{}
	snapSvc := services.NewWorldSnapshotService()
	gs := services.NewGameService(w, logger, snapSvc)

	err := gs.ProcessCommandDTO(dtos.CommandDTO{
		Type:        "SPAWN",
		CharacterID: "hero1",
		Data:        map[string]interface{}{},
	})
	if err != nil {
		t.Fatalf("SPAWN command failed: %v", err)
	}
	if w.Characters["hero1"] == nil {
		t.Fatalf("Expected hero1 to be spawned")
	}

	moveCmd := dtos.CommandDTO{
		Type:        "MOVE",
		CharacterID: "hero1",
		Data: map[string]interface{}{
			"x": 10.0,
			"y": 15.0,
		},
	}
	err = gs.ProcessCommandDTO(moveCmd)
	if err != nil {
		t.Fatalf("MOVE command failed: %v", err)
	}

	x, y := w.Characters["hero1"].Position()
	if x == 0 && y == 0 {
		t.Errorf("hero1 did not move from (0,0)")
	}

	err = gs.ProcessCommandDTO(dtos.CommandDTO{
		Type:        "SPAWN",
		CharacterID: "enemy1",
		Data:        map[string]interface{}{},
	})
	if err != nil {
		t.Fatalf("SPAWN enemy failed: %v", err)
	}
	enemy := w.Characters["enemy1"]
	if enemy == nil {
		t.Fatalf("Enemy not spawned")
	}
	enemyMage := domain.NewMage("enemy1", 10, 15)
	w.Characters["enemy1"] = enemyMage

	w.Characters["hero1"] = domain.NewWarrior("hero1", x, y)

	attackCmd := dtos.CommandDTO{
		Type:        "ATTACK",
		CharacterID: "hero1",
		Data:        map[string]interface{}{},
	}
	err = gs.ProcessCommandDTO(attackCmd)
	if err != nil {
		t.Fatalf("ATTACK command failed: %v", err)
	}

	if len(logger.events) == 0 {
		t.Errorf("No events were logged")
	}
}

func TestGameService_BroadcastState(t *testing.T) {
	w := domain.NewWorld(50, 50)
	logger := &mockLogger{}
	snapSvc := services.NewWorldSnapshotService()
	gs := services.NewGameService(w, logger, snapSvc)

	w.Characters["dead1"] = domain.NewMage("dead1", 10, 10)
	w.Characters["dead1"].TakeDamage(9999, domain.Physical)

	ctx := context.Background()
	gs.BroadcastState(ctx)

	if w.Characters["dead1"].IsDead() {
		t.Errorf("Character dead1 should be respawned, but isDead()=true")
	}
}
