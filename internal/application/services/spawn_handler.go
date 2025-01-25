package services

import (
	"meatgrinder/internal/application/command"
	"meatgrinder/internal/domain"
)

type SpawnHandler struct {
	logger Logger
	world  *domain.World
}

func NewSpawnHandler(world *domain.World, logger Logger) *SpawnHandler {
	return &SpawnHandler{
		world:  world,
		logger: logger,
	}
}

func (h *SpawnHandler) Handle(c command.Command) error {
	if _, ok := h.world.Characters[c.CharacterID]; ok {
		return nil
	}

	h.world.SpawnRandomCharacter(c.CharacterID)
	return nil
}
