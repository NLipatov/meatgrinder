package services

import (
	"fmt"
	"meatgrinder/internal/application/command"
	"meatgrinder/internal/domain"
)

type DisconnectHandler struct {
	world  *domain.World
	logger Logger
}

func NewDisconnectHandler(w *domain.World, l Logger) *DisconnectHandler {
	return &DisconnectHandler{
		world:  w,
		logger: l,
	}
}

func (h *DisconnectHandler) Handle(c command.Command) error {
	delete(h.world.Characters, c.CharacterID)
	h.logger.LogEvent(fmt.Sprintf("%s disconnected", c.CharacterID))
	return nil
}
