package services

import (
	"context"
	"fmt"
	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/domain"
)

type GameService struct {
	world *domain.World
	snap  *WorldSnapshotService
}

func NewGameService(w *domain.World, s *WorldSnapshotService) *GameService {
	return &GameService{world: w, snap: s}
}

func (gs *GameService) ProcessCommandDTO(d dtos.CommandDTO) error {
	cmd, err := MapDTOToCommand(d)
	if err != nil {
		return err
	}
	return gs.ProcessCommand(cmd)
}

func (gs *GameService) ProcessCommand(c Command) error {
	switch c.Type {
	case "SPAWN":
		return gs.handleSpawn(c.CharacterID)
	case "MOVE":
		return gs.handleMove(c.CharacterID, c.Data)
	case "ATTACK":
		return gs.handleAttack(c.CharacterID, c.Data)
	default:
		return fmt.Errorf("unknown cmd %s", c.Type)
	}
}

func (gs *GameService) handleSpawn(id string) error {
	if _, ok := gs.world.Characters[id]; ok {
		return nil
	}
	gs.world.SpawnRandomCharacter(id)
	return nil
}

func (gs *GameService) handleMove(id string, data map[string]interface{}) error {
	ch, ok := gs.world.Characters[id]
	if !ok {
		return nil
	}
	if ch.IsDead() {
		return nil
	}
	dx, _ := data["dx"].(float64)
	dy, _ := data["dy"].(float64)
	ch.MoveStep(dx, dy)
	return nil
}

func (gs *GameService) handleAttack(id string, data map[string]interface{}) error {
	ch, ok := gs.world.Characters[id]
	if !ok {
		return nil
	}
	if ch.IsDead() {
		return nil
	}
	tid, _ := data["target_id"].(string)
	tg, exist := gs.world.Characters[tid]
	if !exist || tg.IsDead() {
		return nil
	}
	ch.Attack([]domain.Character{tg})
	return nil
}

func (gs *GameService) UpdateWorld() {
	gs.world.Update()
}

func (gs *GameService) BroadcastState(ctx context.Context) {
	for _, c := range gs.world.Characters {
		if c.IsDead() && c.State() == domain.StateDying {
			id := c.ID()
			gs.world.SpawnRandomCharacter(id)
		}
	}
}

func (gs *GameService) BuildWorldSnapshot() WorldSnapshot {
	return gs.snap.BuildSnapshot(gs.world)
}
