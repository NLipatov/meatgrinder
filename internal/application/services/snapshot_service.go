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
	Class  string  `json:"class"`
	State  string  `json:"state"`
	Health float64 `json:"health"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Flash  bool    `json:"flash"`
}

func (svc *WorldSnapshotService) BuildSnapshot(w *domain.World) WorldSnapshot {
	var snap WorldSnapshot
	for _, ch := range w.Characters {
		xx, yy := ch.Position()
		cc := "warrior"
		switch ch.(type) {
		case *domain.Mage:
			cc = "mage"
		}
		snap.Characters = append(snap.Characters, CharacterSnapshot{
			ID:     ch.ID(),
			Class:  cc,
			State:  string(ch.State()),
			Health: ch.Health(),
			X:      xx,
			Y:      yy,
			Flash:  ch.FlashRed(),
		})
	}
	return snap
}
