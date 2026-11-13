package ipc

import (
	"bytes"
	"io"
	"log"

	"github.com/jezek/xgbutil"

	"github.com/Ankumeah/NilSky/ipc/commands"
)

func handleCommand(x *xgbutil.XUtil, w io.Writer, data []byte) {
	fields := bytes.Fields(data)
	if len(fields) < 1 {
		w.Write([]byte("ERROR: Empty command"))
		return
	}

	cmd, ok := commands.Commands[string(fields[0])]
	if !ok {
		w.Write([]byte("ERROR: Invalid command"))
		return
	}

	if err := cmd(x, w, fields[1:]...); err != nil {
		log.Printf("Error executing sock command %s: %v\n", string(fields[0]), err)
		w.Write([]byte("ERROR: Check WM logs"))
		return
	}
}
