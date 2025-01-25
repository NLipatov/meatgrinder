package services

import (
	"fmt"
	"math"
	"meatgrinder/internal/application/command"
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

func (h *MoveHandler) Handle(c command.Command) error {
	ch, ok := h.world.Characters[c.CharacterID]
	if !ok {
		return fmt.Errorf("character not found")
	}
	if ch.IsDead() {
		return fmt.Errorf("character is dead")
	}

	cx, cy := ch.Position()

	dxVal, ok := c.Data["dx"].(float64)
	if !ok {
		return fmt.Errorf("invalid dx value")
	}
	dyVal, ok := c.Data["dy"].(float64)
	if !ok {
		return fmt.Errorf("invalid dy value")
	}

	nx := cx + dxVal
	ny := cy + dyVal

	nx = h.clamp(nx, 0, h.world.Width)
	ny = h.clamp(ny, 0, h.world.Height)

	deltaX := nx - cx
	deltaY := ny - cy

	ch.MoveStep(deltaX, deltaY)

	acx, acy := ch.Position()
	distance := math.Hypot(acx-cx, acy-cy)
	if distance > 50 {
		distance = 50
	}

	if distance > 0 {
		h.logMove(ch, distance)
	}

	return nil
}

func (h *MoveHandler) clamp(v, minV, maxV float64) float64 {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func (h *MoveHandler) logMove(character domain.Character, distance float64) {
	h.logger.LogEvent(fmt.Sprintf("%s moved (distance: %v)", character.ID(), distance))
}
