package commands

import (
	"io"
	"strconv"

	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/ewmh"
)

func init() {
	Commands["close_window"] = close_window
}

func close_window(x *xgbutil.XUtil, w io.Writer, args ...[]byte) error {
	if len(args) < 1 {
		w.Write([]byte("ERROR: Window not specified"))
		return nil
	} else if len(args) > 1 {
		w.Write([]byte("ERROR: Expected only one window"))
		return nil
	}

	var winID xproto.Window

	if string(args[0]) == "focus" {
		info, err := xproto.GetInputFocus(x.Conn()).Reply()
		if err != nil {
			return err
		}

		winID = info.Focus
	} else {
		inputPID, err := strconv.Atoi(string(args[0]))
		if err != nil {
			w.Write([]byte("ERROR: Invalid PID"))
			return nil
		}

		wins, err := ewmh.ClientListGet(x)
		if err != nil {
			return err
		}

		found := false
		for _, win := range wins {
			pid, err := ewmh.WmPidGet(x, win)
			if err != nil {
				return err
			}

			if pid == uint(inputPID) {
				winID = win
				found = true
				break
			}
		}

		if !found {
			w.Write([]byte("ERROR: Invalid PID"))
			return nil
		}
	}

	if err := ewmh.CloseWindow(x, winID); err != nil {
		return err
	}

	w.Write([]byte("OK: Window closed"))

	return nil
}
