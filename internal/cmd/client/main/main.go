package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"meatgrinder/internal/application/commands"
	"meatgrinder/internal/cmd/settings"
	"path/filepath"
	"strings"
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

type Fireball struct {
	x, y             float64
	targetX, targetY float64
	img              *ebiten.Image
	timer            float64
}

type Game struct {
	ctx                     context.Context
	cancel                  context.CancelFunc
	client                  *network.Client
	id                      string
	w, h                    int
	mu                      sync.Mutex
	snap                    WorldSnapshot
	prevSnap                WorldSnapshot
	fireballs               []Fireball
	bg                      *ebiten.Image
	mIdle, mRun, mAtk, mDie *ebiten.Image
	wIdle, wRun, wAtk, wDie *ebiten.Image
	fireballImg             *ebiten.Image
	spW, spH                float64
	speed                   float64
}

func NewGame(addr, id string) (*Game, error) {
	ctx, c := context.WithCancel(context.Background())
	g := &Game{
		ctx:       ctx,
		cancel:    c,
		id:        id,
		w:         settings.MapWidth,
		h:         settings.MapHeight,
		speed:     2,
		fireballs: []Fireball{},
	}
	cl := network.NewClient(addr)
	if err := cl.Connect(ctx); err != nil {
		return nil, err
	}
	g.client = cl
	go g.listen()

	spawnCmd := dtos.CommandDTO{
		Type:        commands.SPAWN,
		CharacterID: id,
		Data:        map[string]interface{}{},
	}
	_ = g.client.SendCommand(spawnCmd)

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
	var target *CharacterSnapshot
	var tid string

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.sendMoveCommand(0, -g.speed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.sendMoveCommand(0, g.speed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.sendMoveCommand(-g.speed, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.sendMoveCommand(g.speed, 0)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		tid = g.findCharUnder(float64(mx), float64(my))
		if tid != "" && tid != g.id {
			_ = g.client.SendCommand(dtos.CommandDTO{
				Type:        commands.ATTACK,
				CharacterID: g.id,
				Data:        map[string]interface{}{"target_id": tid},
			})

			g.mu.Lock()
			for _, c := range g.snap.Characters {
				if c.ID == tid {
					target = &c
					break
				}
			}
			var attacker *CharacterSnapshot
			for _, c := range g.snap.Characters {
				if c.ID == g.id {
					attacker = &c
					break
				}
			}
			g.mu.Unlock()

			if target != nil && attacker != nil && strings.ToLower(attacker.Class) == "mage" {
				fb := Fireball{
					x:       attacker.X,
					y:       attacker.Y,
					targetX: target.X,
					targetY: target.Y,
					img:     g.fireballImg,
					timer:   1.0,
				}
				g.mu.Lock()
				g.fireballs = append(g.fireballs, fb)
				g.mu.Unlock()
			}
		}
	}

	g.updateFireballs(1.0 / 60.0)
	return nil
}

func (g *Game) sendMoveCommand(dx, dy float64) {
	_ = g.client.SendCommand(dtos.CommandDTO{
		Type:        commands.MOVE,
		CharacterID: g.id,
		Data: map[string]interface{}{
			"dx": dx,
			"dy": dy,
		},
	})
}

func (g *Game) updateFireballs(dt float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	newFireballs := g.fireballs[:0]
	for i := range g.fireballs {
		fb := &g.fireballs[i]
		fb.timer -= dt
		if fb.timer > 0 {
			dx := fb.targetX - fb.x
			dy := fb.targetY - fb.y
			dist := math.Hypot(dx, dy)
			if dist > 0 {
				speed := 1000.0
				fb.x += (dx / dist) * speed * dt
				fb.y += (dy / dist) * speed * dt
			}
			if fb.x < 0 || fb.x > float64(g.w) || fb.y < 0 || fb.y > float64(g.h) {
				continue
			}
			newFireballs = append(newFireballs, g.fireballs[i])
		}
	}
	g.fireballs = newFireballs
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.bg != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.bg, op)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, fb := range g.fireballs {
		op := &ebiten.DrawImageOptions{}
		scale := 0.1
		op.GeoM.Scale(scale, scale)
		fbWidth, fbHeight := fb.img.Size()
		op.GeoM.Translate(fb.x-float64(fbWidth)*scale/2, fb.y-float64(fbHeight)*scale/2)
		screen.DrawImage(fb.img, op)
	}

	for _, c := range g.snap.Characters {
		img := g.pickSprite(c.Class, c.State)
		if strings.ToLower(c.State) == "dying" {
			img = g.dyingSprite(c.Class)
		}
		op := &ebiten.DrawImageOptions{}
		if c.Flash {
			op.ColorM.Scale(1, 0, 0, 1)
		}
		scale := 0.5
		op.GeoM.Scale(scale, scale)
		charWidth, charHeight := img.Size()
		op.GeoM.Translate(c.X-float64(charWidth)*scale/2, c.Y-float64(charHeight)*scale/2)
		screen.DrawImage(img, op)
	}
}

func (g *Game) pickSprite(class, st string) *ebiten.Image {
	lowerClass := strings.ToLower(class)
	lowerState := strings.ToLower(st)
	if lowerClass == "mage" {
		if lowerState == "running" {
			return g.mRun
		}
		if lowerState == "attacking" {
			return g.mAtk
		}
		return g.mIdle
	}
	if lowerClass == "warrior" {
		if lowerState == "running" {
			return g.wRun
		}
		if lowerState == "attacking" {
			return g.wAtk
		}
		return g.wIdle
	}
	return g.wIdle
}

func (g *Game) dyingSprite(class string) *ebiten.Image {
	lowerClass := strings.ToLower(class)
	if lowerClass == "mage" {
		return g.mDie
	}
	return g.wDie
}

func (g *Game) findCharUnder(mx, my float64) string {
	for _, c := range g.snap.Characters {
		sw := 0.5 * g.spW
		sh := 0.5 * g.spH
		if mx >= c.X-sw/2 && mx <= c.X+sw/2 && my >= c.Y-sh/2 && my <= c.Y+sh/2 {
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

	fireball, _ := loadImg(filepath.Join(dir, "fireball.png"))
	g.fireballImg = fireball

	if m1 != nil {
		ww, hh := m1.Size()
		g.spW, g.spH = float64(ww)*0.5, float64(hh)*0.5
	}
}

func loadImg(path string) (*ebiten.Image, error) {
	f, err := ebitenutil.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func(f ebitenutil.ReadSeekCloser) {
		_ = f.Close()
	}(f)
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
		rand.New(rand.NewSource(time.Now().UnixNano()))
		*id = fmt.Sprintf("player-%04d", rand.Intn(9999))
	}

	g, err := NewGame(*addr, *id)
	if err != nil {
		log.Fatal(err)
	}
	g.LoadAssets(*assetsDir)
	ebiten.SetWindowSize(g.w, g.h)
	ebiten.SetWindowTitle("Meatgrinder")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
