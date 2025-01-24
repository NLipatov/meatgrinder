package domain

import "math/rand"

type World struct {
	Characters map[string]Character
	Width      float64
	Height     float64
}

func NewWorld(w, h float64) *World {
	return &World{
		Characters: make(map[string]Character),
		Width:      w,
		Height:     h,
	}
}

func (wd *World) SpawnRandomCharacter(id string) {
	x := rand.Float64() * wd.Width
	y := rand.Float64() * wd.Height
	if rand.Intn(2) == 0 {
		wd.Characters[id] = NewWarrior(id, x, y)
	} else {
		wd.Characters[id] = NewMage(id, x, y)
	}
}

func (wd *World) Update() {
	dt := 1.0 / 60.0
	for _, c := range wd.Characters {
		c.Update(dt)
	}
}
