package ipc

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/jezek/xgbutil"

	"github.com/Ankumeah/NilSky/meta"
)

func getSocketPath() (string, error) {
	path, ok := os.LookupEnv("XDG_RUNTIME_DIR")
	if !ok {
		return "", errors.New("XDG_RUNTIME_DIR not set")
	}

	return filepath.Join(path, meta.SocketFile), nil
}

func Listen(ctx context.Context, x *xgbutil.XUtil) {
	socketPath, err := getSocketPath()
	if err != nil {
		log.Printf("Error getting socket path: %v\n", err)
		return
	}

	sock, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Printf("Error listing on socket: %v\n", err)
		return
	}
	defer func() {
		sock.Close()
		os.Remove(socketPath)
	}()

	go func() {
		<-ctx.Done()
		sock.Close()
		os.Remove(socketPath)
	}()

	log.Printf("Sock open at %s\n", socketPath)
	for {
		c, err := sock.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			log.Printf("Error accepting sock conn: %v\n", err)
			continue
		}

		go func() {
			defer c.Close()

			data, err := io.ReadAll(c)
			if err != nil {
				log.Printf("Error reading conn data: %v\n", err)
			}

			handleCommand(x, c, data)
		}()
	}
}
