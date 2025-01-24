package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"path/filepath"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/infrastructure/network"
)

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

type Game struct {
	ctx                     context.Context
	cancel                  context.CancelFunc
	client                  *network.Client
	id                      string
	w, h                    int
	mu                      sync.Mutex
	snap                    WorldSnapshot
	bg                      *ebiten.Image
	mIdle, mRun, mAtk, mDie *ebiten.Image
	wIdle, wRun, wAtk, wDie *ebiten.Image
	spW, spH                float64
	speed                   float64
}

func NewGame(addr, id string) (*Game, error) {
	ctx, c := context.WithCancel(context.Background())
	g := &Game{
		ctx:    ctx,
		cancel: c,
		id:     id,
		w:      800,
		h:      600,
		speed:  2,
	}
	cl := network.NewClient(addr)
	if err := cl.Connect(ctx); err != nil {
		return nil, err
	}
	g.client = cl
	go g.listen()

	spawnCmd := dtos.CommandDTO{
		Type:        "SPAWN",
		CharacterID: id,
		Data:        map[string]interface{}{},
	}
	g.client.SendCommand(spawnCmd)

	return g, nil
}

func (g *Game) listen() {
	for data := range g.client.UpdatesChannel() {
		b, err := json.Marshal(data)
		if err != nil {
			continue
		}
		var ws WorldSnapshot
		if err := json.Unmarshal(b, &ws); err != nil {
			continue
		}
		g.mu.Lock()
		g.snap = ws
		g.mu.Unlock()
	}
}

func (g *Game) Layout(int, int) (int, int) {
	return g.w, g.h
}

func (g *Game) Update() error {
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy -= g.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy += g.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx -= g.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx += g.speed
	}
	if math.Abs(dx)+math.Abs(dy) > 0 {
		g.client.SendCommand(dtos.CommandDTO{
			Type:        "MOVE",
			CharacterID: g.id,
			Data: map[string]interface{}{
				"dx": dx,
				"dy": dy,
			},
		})
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		tid := g.findCharUnder(float64(mx), float64(my))
		if tid != "" && tid != g.id {
			g.client.SendCommand(dtos.CommandDTO{
				Type:        "ATTACK",
				CharacterID: g.id,
				Data:        map[string]interface{}{"target_id": tid},
			})
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{40, 40, 80, 255})
	if g.bg != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.bg, op)
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, c := range g.snap.Characters {
		img := g.pickSprite(c.Class, c.State)
		if c.State == "dying" {
			img = g.dyingSprite(c.Class)
		}
		op := &ebiten.DrawImageOptions{}
		if c.Flash {
			op.ColorM.Scale(1, 0, 0, 1)
		}
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(c.X, c.Y)
		screen.DrawImage(img, op)
	}
}

func (g *Game) pickSprite(class, st string) *ebiten.Image {
	if class == "mage" {
		if st == "running" {
			return g.mRun
		}
		if st == "attacking" {
			return g.mAtk
		}
		return g.mIdle
	}
	if class == "warrior" {
		if st == "running" {
			return g.wRun
		}
		if st == "attacking" {
			return g.wAtk
		}
		return g.wIdle
	}
	return g.wIdle
}

func (g *Game) dyingSprite(class string) *ebiten.Image {
	if class == "mage" {
		return g.mDie
	}
	return g.wDie
}

func (g *Game) findCharUnder(mx, my float64) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, c := range g.snap.Characters {
		sw := 0.5 * g.spW
		sh := 0.5 * g.spH
		if mx >= c.X && mx <= c.X+sw && my >= c.Y && my <= c.Y+sh {
			return c.ID
		}
	}
	return ""
}

func (g *Game) Close() error {
	g.cancel()
	g.client.Close()
	return nil
}

func (g *Game) LoadAssets(dir string) {
	bg, _ := loadImg(filepath.Join(dir, "background.png"))
	g.bg = bg

	m1, _ := loadImg(filepath.Join(dir, "mage.png"))
	g.mIdle = m1
	m2, _ := loadImg(filepath.Join(dir, "mage-running.png"))
	g.mRun = m2
	m3, _ := loadImg(filepath.Join(dir, "mage-attacking.png"))
	g.mAtk = m3
	m4, _ := loadImg(filepath.Join(dir, "mage-dying.png"))
	g.mDie = m4

	w1, _ := loadImg(filepath.Join(dir, "warrior.png"))
	g.wIdle = w1
	w2, _ := loadImg(filepath.Join(dir, "warrior-running.png"))
	g.wRun = w2
	w3, _ := loadImg(filepath.Join(dir, "warrior-attacking.png"))
	g.wAtk = w3
	w4, _ := loadImg(filepath.Join(dir, "warrior-dying.png"))
	g.wDie = w4

	if m1 != nil {
		ww, hh := m1.Size()
		g.spW, g.spH = float64(ww), float64(hh)
	}
}

func loadImg(path string) (*ebiten.Image, error) {
	f, err := ebitenutil.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err2 := image.Decode(f)
	if err2 != nil {
		return nil, err2
	}
	return ebiten.NewImageFromImage(img), nil
}

func main() {
	addr := flag.String("addr", "localhost:8080", "server addr")
	id := flag.String("id", "", "character ID (empty => random)")
	assetsDir := flag.String("assets", "assets", "path to assets folder")
	flag.Parse()

	if *id == "" {
		rand.Seed(time.Now().UnixNano())
		*id = fmt.Sprintf("player-%04d", rand.Intn(9999))
	}

	g, err := NewGame(*addr, *id)
	if err != nil {
		log.Fatal(err)
	}
	g.LoadAssets(*assetsDir)
	ebiten.SetWindowSize(g.w, g.h)
	ebiten.SetWindowTitle("Meatgrinder (WASD + Attack)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
