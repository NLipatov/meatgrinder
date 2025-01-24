package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"meatgrinder/internal/application/services"
	"meatgrinder/internal/domain"
	"meatgrinder/internal/infrastructure/network"
	"meatgrinder/internal/infrastructure/persistence"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		fmt.Println("Server shutting down...")
		cancel()
	}()

	world := domain.NewWorld(100, 100)
	logger := persistence.NewFileLogger("game_events.log")
	snapSvc := services.NewWorldSnapshotService()
	gameService := services.NewGameService(world, logger, snapSvc)

	srv := network.NewServer(":8080", gameService)

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				world.Update()
				gameService.BroadcastState(ctx)
			}
		}
	}()

	if err := srv.ListenAndServe(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped.")
}
