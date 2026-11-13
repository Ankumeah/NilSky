package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Ankumeah/NilSky/config"
	"github.com/Ankumeah/NilSky/sock"
	"github.com/Ankumeah/NilSky/tile"
	"github.com/jezek/xgbutil/xevent"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	var wg sync.WaitGroup

	go func() {
		<-ctx.Done()

		log.Println("Received close signal, waiting for goroutines")
		wg.Wait()

		log.Println("Exiting. Bye bye :]")
		os.Exit(0)
	}()

	x, err := GetConn()
	if err != nil {
		log.Fatalf("Error getting connection to x11 server: %v\n", err)
	}

	config, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	if err := tile.Tile(x, config); err != nil {
		log.Fatalf("Error setting up tiling: %v\n", err)
	}

	if err := LoadKeybinds(x, config.Keybinds); err != nil {
		log.Fatalf("Error loading keybinds: %v\n", err)
	}

	wg.Go(func() { sock.Listen(ctx) })

	log.Println("Up and running")
	xevent.Main(x)
}
