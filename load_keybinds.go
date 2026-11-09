package main

import (
	"log"
	"os/exec"

	"github.com/Ankumeah/NilSky/config"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/keybind"
	"github.com/jezek/xgbutil/xevent"
)

func LoadKeybinds(x *xgbutil.XUtil, binds []config.Keybind) error {
	keybind.Initialize(x)

	for _, bind := range binds {
		if err := keybind.KeyPressFun(
			func(*xgbutil.XUtil, xevent.KeyPressEvent) {
				cmd := exec.Command(bind.Cmd[0], bind.Cmd[1:]...)
				if err := cmd.Start(); err != nil {
					log.Printf("Error running keybind: %\n", err)
					return
				}
				go cmd.Wait()
			},
		).Connect(x, x.RootWin(), bind.Bind, true); err != nil {
			return err
		}
	}

	return nil
}
