package services

import (
	"context"
	"fmt"
	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/domain"
)

type GameService struct {
	world         *domain.World
	snap          *WorldSnapshotService
	logger        Logger
	attackHandler Handler
	moveHandler   Handler
	SpawnHandler  Handler
}

func NewGameService(w *domain.World, logger Logger, s *WorldSnapshotService) *GameService {
	return &GameService{
		world:         w,
		logger:        logger,
		snap:          s,
		attackHandler: NewAttackHandler(w, logger),
		moveHandler:   NewMoveHandler(w, logger),
		SpawnHandler:  NewSpawnHandler(w, logger),
	}
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
		return gs.SpawnHandler.Handle(c)
	case "MOVE":
		return gs.moveHandler.Handle(c)
	case "ATTACK":
		return gs.attackHandler.Handle(c)
	default:
		return fmt.Errorf("unknown cmd %s", c.Type)
	}
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
