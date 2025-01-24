package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	_ "image/png"
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
	Class  string  `json:"class"`
	Health float64 `json:"health"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

type Game struct {
	netClient *network.Client
	ctx       context.Context
	cancel    context.CancelFunc

	characterID   string
	width, height int

	mu       sync.Mutex
	snapshot WorldSnapshot

	mageImage    *ebiten.Image
	warriorImage *ebiten.Image

	spriteWidth, spriteHeight float64
}

func NewGame(serverAddr, charID string) (*Game, error) {
	g := &Game{
		characterID: charID,
		width:       800,
		height:      600,
	}
	ctx, cancel := context.WithCancel(context.Background())
	g.ctx = ctx
	g.cancel = cancel

	g.netClient = network.NewClient(serverAddr)
	if err := g.netClient.Connect(ctx); err != nil {
		return nil, fmt.Errorf("connect failed: %w", err)
	}
	spawn := dtos.CommandDTO{Type: "SPAWN", CharacterID: charID}
	_ = g.netClient.SendCommand(spawn)

	go g.listenUpdates()

	return g, nil
}

func (g *Game) listenUpdates() {
	for data := range g.netClient.UpdatesChannel() {
		raw, err := json.Marshal(data)
		if err != nil {
			log.Println("Marshal error:", err)
			continue
		}
		var snap WorldSnapshot
		if err := json.Unmarshal(raw, &snap); err != nil {
			log.Println("Unmarshal error:", err)
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
	speed := 8.0
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
		mx, my := ebiten.CursorPosition()
		clickedID := g.findCharacterUnderMouse(float64(mx), float64(my))
		if clickedID != "" {
			cmd := dtos.CommandDTO{
				Type:        "ATTACK",
				CharacterID: g.characterID,
				Data:        map[string]interface{}{"target_id": clickedID},
			}
			_ = g.netClient.SendCommand(cmd)
			log.Printf("ATTACK on %s", clickedID)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 45, G: 45, B: 45, A: 255})

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, ch := range g.snapshot.Characters {
		var img *ebiten.Image
		switch ch.Class {
		case "warrior":
			img = g.warriorImage
		case "mage":
			img = g.mageImage
		default:
			img = g.warriorImage
		}
		if img == nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.1, 0.1)
		op.GeoM.Translate(ch.X, ch.Y)
		screen.DrawImage(img, op)
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

func (g *Game) findCharacterUnderMouse(mx, my float64) string {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, c := range g.snapshot.Characters {
		x1, y1 := c.X, c.Y
		x2, y2 := x1+g.spriteWidth, y1+g.spriteHeight
		if mx >= x1 && mx <= x2 && my >= y1 && my <= y2 {
			return c.ID
		}
	}
	return ""
}

func (g *Game) Close() error {
	g.cancel()
	g.netClient.Close()
	return nil
}

func main() {
	serverAddr := flag.String("addr", "localhost:8080", "Server address")
	charID := flag.String("id", "player123", "Character ID")
	flag.Parse()

	g, err := NewGame(*serverAddr, *charID)
	if err != nil {
		log.Fatal(err)
	}

	mageImg, err := loadImage("assets/mage.png")
	if err != nil {
		log.Println("failed to load mage.png:", err)
	}
	g.mageImage = mageImg

	wImg, err := loadImage("assets/warrior.png")
	if err != nil {
		log.Println("failed to load warrior.png:", err)
	}
	g.warriorImage = wImg

	if mageImg != nil {
		w, h := mageImg.Size()
		g.spriteWidth, g.spriteHeight = float64(w), float64(h)
	}

	ebiten.SetWindowSize(g.width, g.height)
	ebiten.SetWindowTitle("Meatgrinder Ebitengine Client (Sprites)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func loadImage(path string) (*ebiten.Image, error) {
	f, err := ebitenutil.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}
