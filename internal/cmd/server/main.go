package main

import (
	"context"
	"log"
	"meatgrinder/internal/application/services"
	"meatgrinder/internal/domain"
	"meatgrinder/internal/infrastructure/network"
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
	w := domain.NewWorld(100, 100)
	svc := services.NewWorldSnapshotService()
	gs := services.NewGameService(w, svc)
	srv := network.NewServer(":8080", gs)

	go func() {
		t := time.NewTicker(time.Second / 60)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				gs.UpdateWorld()
				gs.BroadcastState(ctx)
			}
		}
	}()

	if err := srv.ListenAndServe(ctx); err != nil {
		log.Fatal(err)
	}
}
