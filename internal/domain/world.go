package domain

import "math/rand"

type World struct {
	Characters map[string]Character
	Width      float64
	Height     float64
}

func NewWorld(width, height float64) *World {
	return &World{
		Characters: make(map[string]Character),
		Width:      width,
		Height:     height,
	}
}

func (w *World) SpawnRandomCharacter(id string) {
	x := rand.Float64() * w.Width
	y := rand.Float64() * w.Height

	if rand.Intn(2) == 0 {
		w.Characters[id] = NewWarrior(id, x, y)
	} else {
		w.Characters[id] = NewMage(id, x, y)
	}
}

func (w *World) Update() {
	dt := 1.0 / 20.0
	for _, ch := range w.Characters {
		ch.Update(dt)
	}
}
