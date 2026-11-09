package main

import (
	"log"

	"github.com/Ankumeah/NilSky/config"
	"github.com/jezek/xgbutil/xevent"
)

func main() {
	x, err := GetConn()
	if err != nil {
		log.Fatalf("Error getting connection to x11 server: %v\n", err)
	}

	config, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	if err := LoadKeybinds(x, config.Keybinds); err != nil {
		log.Fatalf("Error loading keybinds: %v\n", err)
	}

	log.Println("Up and running")
	xevent.Main(x)
}
