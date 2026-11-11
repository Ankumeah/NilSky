package tile

import (
	"sync"

	"github.com/Ankumeah/NilSky/layouts"
)

type tileInfo struct {
	workspaces map[uint]layouts.Layout
	currentWS  uint
	mu         sync.Mutex
}

type window struct {
	size vec2
	pos  vec2
}

type vec2 struct {
	x uint
	y uint
}
