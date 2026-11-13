package commands

import (
	"io"

	"github.com/jezek/xgbutil"
)

func init() {
	Commands["ping"] = ping
}

func ping(x *xgbutil.XUtil, w io.Writer, args ...[]byte) error {
	w.Write([]byte("OK: PONG"))

	return nil
}
