package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"log"
	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/infrastructure/network"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type WorldSnapshot struct {
	Characters []CharacterSnapshot `json:"characters"`
}

type CharacterSnapshot struct {
	ID     string  `json:"id"`
	Health float64 `json:"health"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

type Game struct {
	netClient *network.Client
	ctx       context.Context
	cancel    context.CancelFunc

	characterID string

	mu       sync.Mutex
	snapshot WorldSnapshot

	width, height int
}

func NewGame(serverAddr string, characterID string) *Game {
	ctx, cancel := context.WithCancel(context.Background())
	return &Game{
		characterID: characterID,
		ctx:         ctx,
		cancel:      cancel,
		width:       800,
		height:      600,
	}
}

func (g *Game) ConnectServer(serverAddr string) error {
	g.netClient = network.NewClient(serverAddr)
	if err := g.netClient.Connect(g.ctx); err != nil {
		return fmt.Errorf("connect server: %w", err)
	}

	go g.listenServerUpdates()

	spawnCmd := dtos.CommandDTO{
		Type:        "SPAWN",
		CharacterID: g.characterID,
		Data:        map[string]interface{}{},
	}
	_ = g.netClient.SendCommand(spawnCmd)

	log.Println("Connected to server at", serverAddr, "as", g.characterID)
	return nil
}

func (g *Game) listenServerUpdates() {
	defer log.Println("listenServerUpdates finished.")
	for data := range g.netClient.UpdatesChannel() {
		raw, err := json.Marshal(data)
		if err != nil {
			log.Println("marshal error:", err)
			continue
		}
		var snap WorldSnapshot
		if err := json.Unmarshal(raw, &snap); err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		g.mu.Lock()
		g.snapshot = snap
		g.mu.Unlock()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

func (g *Game) Update() error {
	speed := 2.0
	moveX, moveY := 0.0, 0.0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		moveY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		moveY += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		moveX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		moveX += speed
	}

	if moveX != 0 || moveY != 0 {
		x, y := g.getLocalPos()
		cmd := dtos.CommandDTO{
			Type:        "MOVE",
			CharacterID: g.characterID,
			Data: map[string]interface{}{
				"x": x + moveX,
				"y": y + moveY,
			},
		}
		_ = g.netClient.SendCommand(cmd)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cmd := dtos.CommandDTO{
			Type:        "ATTACK",
			CharacterID: g.characterID,
			Data: map[string]interface{}{
				"target_id": g.characterID,
			},
		}
		_ = g.netClient.SendCommand(cmd)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 45, G: 45, B: 45, A: 255})

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, ch := range g.snapshot.Characters {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(ch.X), float64(ch.Y))
		rectImg := ebiten.NewImage(10, 10)
		rectImg.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})

		screen.DrawImage(rectImg, op)
	}
}

func (g *Game) getLocalPos() (float64, float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, c := range g.snapshot.Characters {
		if c.ID == g.characterID {
			return c.X, c.Y
		}
	}
	return 0, 0
}

func (g *Game) Close() error {
	g.cancel()
	g.netClient.Close()
	return nil
}

func main() {
	serverAddr := flag.String("addr", "localhost:8080", "Server address")
	characterID := flag.String("id", "player123", "Character ID")
	flag.Parse()

	game := NewGame(*serverAddr, *characterID)
	if err := game.ConnectServer(*serverAddr); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	ebiten.SetWindowSize(game.width, game.height)
	ebiten.SetWindowTitle("Meatgrinder Ebitengine Client")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Ebiten run error: %v", err)
	}
}
