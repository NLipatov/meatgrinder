package services

import (
	"fmt"
	"math"
	"meatgrinder/internal/domain"
)

type MoveHandler struct {
	logger Logger
	world  *domain.World
}

func NewMoveHandler(world *domain.World, logger Logger) *MoveHandler {
	return &MoveHandler{
		logger: logger,
		world:  world,
	}
}

func (h *MoveHandler) Handle(c Command) error {
	ch, ok := h.world.Characters[c.CharacterID]
	if !ok {
		return nil
	}
	if ch.IsDead() {
		return nil
	}
	dx, _ := c.Data["dx"].(float64)
	dy, _ := c.Data["dy"].(float64)
	ch.MoveStep(dx, dy)

	acx, acy := ch.Position()
	math.Hypot(dx-acx, dy-acy)
	h.logMove(ch, math.Hypot(dx-acx, dy-acy))

	return nil
}

func (h *MoveHandler) logMove(character domain.Character, distance float64) {
	h.logger.LogEvent(fmt.Sprintf("%s moved (distance: %v)", character.ID(), distance))
}
