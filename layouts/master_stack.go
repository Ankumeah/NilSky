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

	"github.com/Ankumeah/NilSky/config"
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/keybind"
	"github.com/jezek/xgbutil/mousebind"
	"github.com/jezek/xgbutil/xevent"
	"github.com/jezek/xgbutil/xwindow"
)

type masterStackWindow struct {
	id     xproto.Window
	x      uint
	y      uint
	width  uint
	height uint
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

	l.height = l.Screen.Size.Heigh - (2 * decor.GapsOut)
	l.width = l.Screen.Size.Width - (2 * decor.GapsOut)

	l.masterWidth = (l.Mfact * (l.width - decor.GapsIn)) / 100
	l.stackWidth = l.width - l.masterWidth - decor.GapsIn

	if l.MLeft {
		l.masterX = decor.GapsOut
		l.stackX = decor.GapsOut + l.masterWidth + decor.GapsIn
	} else {
		l.masterX = decor.GapsOut + l.stackWidth + decor.GapsIn
		l.masterX = decor.GapsOut
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

	if err := xproto.MapWindowChecked(x.Conn(), e.Window).Check(); err != nil {
		return fmt.Errorf("Error mapping window: %v\n", err)
	}

	xevent.DestroyNotifyFun(
		func(x *xgbutil.XUtil, e xevent.DestroyNotifyEvent) {
			if err := l.destroy(x, e); err != nil {
				log.Printf("Error removing window: %v\n", err)
			}
		},
	).Connect(x, e.Window)

	if err := xwindow.New(x, e.Window).Listen(xproto.EventMaskEnterWindow); err != nil {
		return fmt.Errorf("Error listing for enter window: %w\n", err)
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

	return nil
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

	if err := l.redraw(x); err != nil {
		return err
	}

	xevent.Detach(x, e.Window)
	keybind.Detach(x, e.Window)
	mousebind.Detach(x, e.Window)

	return nil
}

func (*MasterStack) focus(x *xgbutil.XUtil, id xproto.Window) error {
	return xproto.SetInputFocusChecked(
		x.Conn(),
		xproto.InputFocusPointerRoot,
		id,
		xproto.TimeCurrentTime,
	).Check()
}

func (l *MasterStack) redraw(x *xgbutil.XUtil) error {
	if l.master == nil {
		return nil
	}

	if err := xproto.ConfigureWindowChecked(
		x.Conn(), l.master.id,
		xproto.ConfigWindowX|
			xproto.ConfigWindowY|
			xproto.ConfigWindowWidth|
			xproto.ConfigWindowHeight,
		[]uint32{
			uint32(l.master.x),
			uint32(l.master.y),
			uint32(l.master.width),
			uint32(l.master.height),
		},
	).Check(); err != nil {
		return fmt.Errorf("Error resizing window: %v\n", err)
	}

	for _, win := range l.stack {
		if err := xproto.ConfigureWindowChecked(
			x.Conn(), win.id,
			xproto.ConfigWindowX|
				xproto.ConfigWindowY|
				xproto.ConfigWindowWidth|
				xproto.ConfigWindowHeight,
			[]uint32{
				uint32(win.x),
				uint32(win.y),
				uint32(win.width),
				uint32(win.height),
			},
		).Check(); err != nil {
			return fmt.Errorf("Error resizing window: %v\n", err)
		}
	}

	return nil
}
