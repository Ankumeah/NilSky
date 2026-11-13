package commands

import (
	"io"

	"github.com/jezek/xgbutil"
)

var Commands = map[string]func(*xgbutil.XUtil, io.Writer, ...[]byte) error{}
