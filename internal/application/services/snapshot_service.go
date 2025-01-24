package services

import (
	"meatgrinder/internal/domain"
)

type WorldSnapshotService struct{}

func NewWorldSnapshotService() *WorldSnapshotService {
	return &WorldSnapshotService{}
}

type WorldSnapshot struct {
	Characters []CharacterSnapshot `json:"characters"`
}

type CharacterSnapshot struct {
	ID     string  `json:"id"`
	Health float64 `json:"health"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

func (svc *WorldSnapshotService) BuildSnapshot(world *domain.World) WorldSnapshot {
	var snap WorldSnapshot
	for _, ch := range world.Characters {
		x, y := ch.Position()
		snap.Characters = append(snap.Characters, CharacterSnapshot{
			ID:     ch.ID(),
			Health: ch.Health(),
			X:      x,
			Y:      y,
		})
	}
	return snap
}
