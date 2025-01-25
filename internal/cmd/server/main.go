package main

import (
	"context"
	"log"
	"meatgrinder/internal/application/services"
	"meatgrinder/internal/cmd/settings"
	"meatgrinder/internal/domain"
	"meatgrinder/internal/infrastructure/network"
	"meatgrinder/internal/infrastructure/persistence"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		cancel()
	}()
	w := domain.NewWorld(settings.MapHeight, settings.MapWidth)
	svc := services.NewWorldSnapshotService()
	l := persistence.NewFileLogger("game_events.log")
	gs := services.NewGameService(w, l, svc)
	srv := network.NewServer(":8080", gs)

	go func() {
		t := time.NewTicker(time.Second / 120)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				gs.UpdateWorld()
				gs.BroadcastState()
			}
		}
	}()

	if err := srv.ListenAndServe(ctx); err != nil {
		log.Fatal(err)
	}
}
