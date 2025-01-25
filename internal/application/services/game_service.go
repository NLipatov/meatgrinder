package services

import (
	"fmt"
	"meatgrinder/internal/application/commands"
	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/domain"
)

type GameService struct {
	world             *domain.World
	snap              *WorldSnapshotService
	logger            Logger
	attackHandler     Handler
	moveHandler       Handler
	spawnHandler      Handler
	disconnectHandler Handler
}

func NewGameService(w *domain.World, logger Logger, s *WorldSnapshotService) *GameService {
	return &GameService{
		world:             w,
		logger:            logger,
		snap:              s,
		attackHandler:     NewAttackHandler(w, logger),
		moveHandler:       NewMoveHandler(w, logger),
		spawnHandler:      NewSpawnHandler(w, logger),
		disconnectHandler: NewDisconnectHandler(w, logger),
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
	case commands.SPAWN:
		return gs.spawnHandler.Handle(c)
	case commands.MOVE:
		return gs.moveHandler.Handle(c)
	case commands.ATTACK:
		return gs.attackHandler.Handle(c)
	case commands.DISCONNECT:
		return gs.disconnectHandler.Handle(c)

	default:
		return fmt.Errorf("unknown cmd %v", c.Type)
	}
}

func (gs *GameService) UpdateWorld() {
	gs.world.Update()
}

func (gs *GameService) BroadcastState() {
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
