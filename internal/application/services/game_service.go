package services

import (
	"context"
	"fmt"
	"math"
	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/domain"
)

type GameService struct {
	world  *domain.World
	logger ILogger
}

func NewGameService(world *domain.World, logger ILogger) *GameService {
	gs := &GameService{
		world:  world,
		logger: logger,
	}
	return gs
}

func (gs *GameService) BroadcastState(ctx context.Context) {
	for _, ch := range gs.world.Characters {
		select {
		case <-ctx.Done():
			return
		default:

			if ch.IsDead() {
				gs.logger.LogEvent(fmt.Sprintf("Character %s is dead", ch.ID()))

				// respawn
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
	ch, ok := gs.world.Characters[cmd.CharacterID]
	if !ok && cmd.Type != "SPAWN" {
		return fmt.Errorf("character %s not found", cmd.CharacterID)
	}

	switch cmd.Type {
	case "MOVE":
		return gs.handleMove(ch, cmd.Data)
	case "ATTACK":
		return gs.handleAttack(ch, cmd.Data)
	case "SPAWN":
		return gs.handleSpawn(cmd.CharacterID, cmd.Data)
	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type)
	}
}

func (gs *GameService) handleMove(ch domain.Character, data map[string]interface{}) error {
	xVal, xOK := data["x"].(float64)
	yVal, yOK := data["y"].(float64)
	if !xOK || !yOK {
		return fmt.Errorf("invalid MOVE data, expect floats x, y")
	}
	ch.MoveTo(xVal, yVal)
	gs.logger.LogEvent(fmt.Sprintf("Character %s moved to (%.2f, %.2f)", ch.ID(), xVal, yVal))
	return nil
}

func (gs *GameService) handleAttack(ch domain.Character, data map[string]interface{}) error {
	if ch.IsDead() {
		return fmt.Errorf("attacker %s is dead", ch.ID())
	}

	attacker := gs.world.Characters[ch.ID()]
	switch attacker.DamageType() {

	case domain.Magical:
		targetID, ok := data["target_id"].(string)
		if !ok || targetID == "" {
			return fmt.Errorf("mage attack requires 'target_id' in data")
		}
		target, exists := gs.world.Characters[targetID]
		if !exists || target.IsDead() {
			return fmt.Errorf("target %s not found or dead", targetID)
		}
		if distance(attacker, target) > 8.0 {
			return fmt.Errorf("target %s is too far for mage's attack", targetID)
		}
		target.TakeDamage(attacker.AttackPower(), domain.Magical)

		if w, isWarrior := target.(*domain.Warrior); isWarrior {
			w.ApplySlow(0.5, 3.0)
		}

		gs.logger.LogEvent(fmt.Sprintf("Mage %s attacked %s with magic", attacker.ID(), targetID))

	case domain.Physical:
		for _, t := range gs.getAllExcept(attacker.ID()) {
			if t.IsDead() {
				continue
			}
			if distance(attacker, t) <= attacker.AttackRadius() {
				t.TakeDamage(attacker.AttackPower(), domain.Physical)
			}
		}

		gs.logger.LogEvent(fmt.Sprintf("Warrior %s did an AoE melee attack", attacker.ID()))

	default:
		return fmt.Errorf("character %s cannot attack (invalid damage type)", ch.ID())
	}
	return nil
}

func (gs *GameService) handleSpawn(charID string, data map[string]interface{}) error {
	if _, exists := gs.world.Characters[charID]; exists {
		return fmt.Errorf("character %s already exists in the world", charID)
	}
	gs.world.SpawnRandomCharacter(charID)
	gs.logger.LogEvent(fmt.Sprintf("Character %s spawned", charID))
	return nil
}

func (gs *GameService) getAllExcept(id string) []domain.Character {
	var others []domain.Character
	for cid, ch := range gs.world.Characters {
		if cid != id && !ch.IsDead() {
			others = append(others, ch)
		}
	}
	return others
}

func distance(a, b domain.Character) float64 {
	ax, ay := a.Position()
	bx, by := b.Position()
	dx := bx - ax
	dy := by - ay
	return math.Hypot(dx, dy)
}
