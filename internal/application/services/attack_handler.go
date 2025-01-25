package services

import (
	"fmt"
	"math"
	"meatgrinder/internal/application/command"
	"meatgrinder/internal/domain"
)

type AttackHandler struct {
	world  *domain.World
	logger Logger
}

func NewAttackHandler(world *domain.World, logger Logger) *AttackHandler {
	return &AttackHandler{
		world:  world,
		logger: logger,
	}
}

func (h *AttackHandler) Handle(c command.Command) error {
	attacker, ok := h.world.Characters[c.CharacterID]
	if !ok {
		return nil
	}
	if attacker.IsDead() {
		return nil
	}
	tid, _ := c.Data["target_id"].(string)
	target, exist := h.world.Characters[tid]
	if !exist || target.IsDead() {
		return nil
	}

	if h.getDistance(attacker, target) > attacker.AttackRadius() {
		return nil
	}

	attacker.Attack([]domain.Character{target})
	h.logAttack(attacker, target)

	return nil
}

func (h *AttackHandler) getDistance(a, t domain.Character) float64 {
	acx, acy := a.Position()
	tcx, tcy := t.Position()
	return math.Hypot(tcx-acx, tcy-acy)
}

func (h *AttackHandler) logAttack(attacker, target domain.Character) {
	switch attacker.DamageType() {
	case domain.Physical:
		h.logger.LogEvent(fmt.Sprintf("Warrior %s attacked %s with sword", attacker.ID(), target.ID()))
	case domain.Magical:
		h.logger.LogEvent(fmt.Sprintf("Mage %s attacked %s with magic", attacker.ID(), target.ID()))
	default:
		h.logger.LogEvent(fmt.Sprintf("%s dealt damage to %s with %s", attacker.ID(), target.ID(), target.DamageType()))
	}
}
