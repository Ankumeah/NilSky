package layouts

import (
	"github.com/Ankumeah/NilSky/config"
)

var Layouts = map[string]func([]byte, config.Decorations) (Layout, error){}
