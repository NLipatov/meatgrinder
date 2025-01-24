package services

import (
	"context"
	"fmt"
	"math"
	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/domain"
)

type GameService struct {
	world           *domain.World
	logger          ILogger
	snapshotService *WorldSnapshotService
}

func NewGameService(
	world *domain.World,
	logger ILogger,
	snapSvc *WorldSnapshotService,
) *GameService {
	return &GameService{
		world:           world,
		logger:          logger,
		snapshotService: snapSvc,
	}
}

func (gs *GameService) BroadcastState(ctx context.Context) {
	for _, ch := range gs.world.Characters {
		select {
		case <-ctx.Done():
			return
		default:
			if ch.IsDead() {
				gs.logger.LogEvent(fmt.Sprintf("Character %s is dead", ch.ID()))
				id := ch.ID()
				gs.world.SpawnRandomCharacter(id)
			}
		}
	}
}

func (gs *GameService) ProcessCommandDTO(dto dtos.CommandDTO) error {
	cmd, err := MapDTOToCommand(dto)
	if err != nil {
		return err
	}
	return gs.ProcessCommand(cmd)
}

func (gs *GameService) ProcessCommand(cmd Command) error {
	switch cmd.Type {
	case "SPAWN":
		return gs.handleSpawn(cmd.CharacterID)

	case "MOVE", "ATTACK":
		ch, exists := gs.world.Characters[cmd.CharacterID]
		if !exists {
			return fmt.Errorf("character %s not found", cmd.CharacterID)
		}
		if ch.IsDead() {
			return fmt.Errorf("character %s is dead and cannot act", ch.ID())
		}

		switch cmd.Type {
		case "MOVE":
			return gs.handleMove(ch, cmd.Data)
		case "ATTACK":
			return gs.handleAttack(ch, cmd.Data)
		}

	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type)
	}
	return nil
}

func (gs *GameService) handleSpawn(charID string) error {
	if _, exists := gs.world.Characters[charID]; exists {
		return fmt.Errorf("character %s already exists in the world", charID)
	}
	gs.world.SpawnRandomCharacter(charID)
	gs.logger.LogEvent(fmt.Sprintf("Character %s spawned", charID))
	return nil
}

func (gs *GameService) handleMove(ch domain.Character, data map[string]interface{}) error {
	xVal, xOK := data["x"].(float64)
	yVal, yOK := data["y"].(float64)
	if !xOK || !yOK {
		return fmt.Errorf("invalid MOVE data: need floats x,y")
	}
	ch.MoveTo(xVal, yVal)
	gs.logger.LogEvent(fmt.Sprintf("Character %s moved to (%.1f, %.1f)", ch.ID(), xVal, yVal))
	return nil
}

func (gs *GameService) handleAttack(ch domain.Character, data map[string]interface{}) error {
	switch ch.DamageType() {
	case domain.Magical:
		targetID, ok := data["target_id"].(string)
		if !ok || targetID == "" {
			return fmt.Errorf("mage attack requires 'target_id'")
		}
		target, exists := gs.world.Characters[targetID]
		if !exists || target.IsDead() {
			return fmt.Errorf("target %s not found or dead", targetID)
		}
		dist := gs.distance(ch, target)
		if dist > ch.AttackRadius() {
			return fmt.Errorf("target %s is too far (%.1f) for radius %.1f", targetID, dist, ch.AttackRadius())
		}
		ch.Attack([]domain.Character{target})
		if w, isWarrior := target.(*domain.Warrior); isWarrior {
			w.ApplySlow(0.5, 3.0)
		}
		gs.logger.LogEvent(fmt.Sprintf("Mage %s attacked %s with magic", ch.ID(), targetID))

	case domain.Physical:
		var victims []domain.Character
		for cid, c := range gs.world.Characters {
			if cid == ch.ID() || c.IsDead() {
				continue
			}
			if gs.distance(ch, c) <= ch.AttackRadius() {
				victims = append(victims, c)
			}
		}
		ch.Attack(victims)
		gs.logger.LogEvent(fmt.Sprintf("Warrior %s did an AoE melee attack", ch.ID()))

	default:
		return fmt.Errorf("unknown damage type for character %s", ch.ID())
	}
	return nil
}

func (gs *GameService) distance(a, b domain.Character) float64 {
	ax, ay := a.Position()
	bx, by := b.Position()
	return math.Hypot(bx-ax, by-ay)
}

func (gs *GameService) BuildWorldSnapshot() WorldSnapshot {
	return gs.snapshotService.BuildSnapshot(gs.world)
}
