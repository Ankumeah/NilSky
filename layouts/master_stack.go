// [WARNING]
//
// THIS IS PROBABLY ONE OF THE
// UGLIEST FILES I HAVE EVER WRITTEN
// AND IV WROTE A LOT OF UGLY FILES
//
// YOU HAVE BEEN WARNED

package layouts

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/keybind"
	"github.com/jezek/xgbutil/mousebind"
	"github.com/jezek/xgbutil/xevent"
	"github.com/jezek/xgbutil/xprop"
	"github.com/jezek/xgbutil/xwindow"

	"github.com/Ankumeah/NilSky/config"
)

type masterStackWindow struct {
	id      xproto.Window
	frameID xproto.Window
	x       uint
	y       uint
	width   uint
	height  uint
}

type MasterStack struct {
	Mfact  uint `json:"mfact"`
	MLeft  bool `json:"mleft"`
	Screen struct {
		Size struct {
			Width uint `json:"width"`
			Heigh uint `json:"height"`
		} `json:"size"`
	} `json:"screen"`

	decor config.Decorations

	width  uint
	height uint

	masterWidth uint
	stackWidth  uint

	borderColor uint

	masterX uint
	stackX  uint

	master *masterStackWindow
	stack  []*masterStackWindow
}

func init() {
	Layouts["master_stack"] = GetMasterStack
}

func GetMasterStack(config []byte, decor config.Decorations) (Layout, error) {
	l := &MasterStack{
		decor: decor,
	}

	if err := json.Unmarshal(config, &l); err != nil {
		return l, err
	}

	if l.Mfact > 100 {
		return l, errors.New("master stack mfact must be in range 1 to 100")
	} else if l.Mfact == 0 {
		return l, errors.New("master stack mfact must not be 0")
	}

	if !strings.HasPrefix(l.decor.BorderColor, "#") {
		return l, errors.New("border_color must start with #")
	}
	borderColor, err := strconv.ParseUint(strings.TrimPrefix(l.decor.BorderColor, "#"), 16, 0)
	if err != nil {
		return l, errors.New("border_color must be a valid hex number")
	}
	l.borderColor = uint(borderColor)

	l.height = l.Screen.Size.Heigh - (2 * decor.GapsOut)
	l.width = l.Screen.Size.Width - (2 * decor.GapsOut)

	l.masterWidth = (l.Mfact * (l.width - decor.GapsIn)) / 100
	l.stackWidth = l.width - l.masterWidth - decor.GapsIn

	if l.MLeft {
		l.masterX = decor.GapsOut
		l.stackX = decor.GapsOut + l.masterWidth + decor.GapsIn
	} else {
		l.masterX = decor.GapsOut + l.stackWidth + decor.GapsIn
		l.stackX = decor.GapsOut
	}

	return l, nil
}

func (l *MasterStack) Map(x *xgbutil.XUtil, e xevent.MapRequestEvent) error {
	w := &masterStackWindow{
		id: e.Window,
	}

	if l.master == nil {
		l.master = w

		w.x = l.decor.GapsOut
		w.y = l.decor.GapsOut
		w.width = l.width
		w.height = l.height
	} else if len(l.stack) == 0 {
		l.stack = append(l.stack, w)

		w.x = l.stackX
		w.y = l.decor.GapsOut
		w.width = l.stackWidth
		w.height = l.height

		l.master.width = l.masterWidth
		l.master.x = l.masterX
	} else {
		l.stack = append(l.stack, w)

		w.x = l.stackX
		w.width = l.stackWidth

		height := (l.height - (uint(len(l.stack)-1) * l.decor.GapsIn)) / uint(len(l.stack))

		for i, win := range l.stack {
			win.y = l.decor.GapsOut + (uint(i) * (height + l.decor.GapsIn))
			win.height = height
		}
	}

	if err := l.redraw(x); err != nil {
		return err
	}

	if err := xproto.MapWindowChecked(
		x.Conn(),
		e.Window,
	).Check(); err != nil {
		return fmt.Errorf("Error mapping window: %v\n", err)
	}

	xevent.DestroyNotifyFun(
		func(x *xgbutil.XUtil, e xevent.DestroyNotifyEvent) {
			if err := l.destroy(x, e); err != nil {
				log.Printf("Error removing window: %v\n", err)
			}
		},
	).Connect(x, e.Window)

	if err := xwindow.New(x, e.Window).Listen(
		xproto.EventMaskEnterWindow,
	); err != nil {
		return fmt.Errorf("Error listing for EventMaskEnterWindow: %w\n", err)
	}

	xevent.EnterNotifyFun(
		func(x *xgbutil.XUtil, e xevent.EnterNotifyEvent) {
			if e.Detail == xproto.NotifyDetailInferior {
				return
			}

			if err := l.focus(x, e.Event); err != nil {
				log.Printf("Error focusing window: %v\n", err)
			}
		},
	).Connect(x, e.Window)

	xevent.ClientMessageFun(
		func(x *xgbutil.XUtil, e xevent.ClientMessageEvent) {
			netCloseWindow, err := xprop.Atm(x, "_NET_CLOSE_WINDOW")
			if err != nil {
				log.Printf("Error getting _NET_CLOSE_WINDOW atom: %v\n", err)
				return
			}

			if e.Type == netCloseWindow {
				handleNetCloseWindow(x, e)
			}
		},
	).Connect(x, e.Window)

	return nil
}

func handleNetCloseWindow(x *xgbutil.XUtil, e xevent.ClientMessageEvent) {
	wmProtocols, err := xprop.Atm(x, "WM_PROTOCOLS")
	if err != nil {
		log.Printf("Error getting WM_PROTOCOLS atom: %v\n", err)
		return
	}
	wmDelete, err := xprop.Atm(x, "WM_DELETE_WINDOW")
	if err != nil {
		log.Printf("Error getting WM_DELETE_WINDOW atom: %v\n", err)
		return
	}

	if err := xproto.SendEventChecked(
		x.Conn(), false, e.Window, 0,
		string(
			xproto.ClientMessageEvent{
				Format: 32,
				Window: e.Window,
				Type:   wmProtocols,
				Data: xproto.ClientMessageDataUnionData32New(
					[]uint32{uint32(wmDelete), 0, 0, 0, 0},
				),
			}.Bytes(),
		),
	).Check(); err != nil {
		log.Printf("Error sending WM_DELETE_WINDOW: %v\n", err)
		return
	}
}

func (l *MasterStack) destroy(x *xgbutil.XUtil, e xevent.DestroyNotifyEvent) error {
	if len(l.stack) == 0 {
		l.master = nil
	} else if len(l.stack) == 1 {
		if e.Window == l.master.id {
			l.master = l.stack[0]
		}
		l.stack = []*masterStackWindow{}

		l.master.x = l.decor.GapsOut
		l.master.y = l.decor.GapsOut

		l.master.width = l.width
		l.master.height = l.height
	} else if e.Window == l.master.id {
		l.master.id = l.stack[0].id
		l.master.frameID = l.stack[0].frameID
		l.stack = l.stack[1:]

		height := (l.height - (uint(len(l.stack)-1) * l.decor.GapsIn)) / uint(len(l.stack))

		for i, win := range l.stack {
			win.y = l.decor.GapsOut + (uint(i) * (height + l.decor.GapsIn))
			win.height = height
		}

	} else {
		for i, win := range l.stack {
			if win.id == e.Window {
				l.stack = append(l.stack[:i], l.stack[i+1:]...)
			}
		}

		height := (l.height - (uint(len(l.stack)-1) * l.decor.GapsIn)) / uint(len(l.stack))

		for i, win := range l.stack {
			win.height = height
			win.y = l.decor.GapsOut + (uint(i) * (height + l.decor.GapsIn))
		}
	}

	if err := xproto.DestroyWindowChecked(
		x.Conn(),
		e.Event,
	).Check(); err != nil {
		return err
	}

	xevent.Detach(x, e.Window)
	keybind.Detach(x, e.Window)
	mousebind.Detach(x, e.Window)

	if err := l.redraw(x); err != nil {
		return err
	}

	return nil
}
