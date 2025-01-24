package services

import (
	"context"
	"fmt"
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
