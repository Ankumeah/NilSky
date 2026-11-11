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
	id xproto.Window
}

type MasterStack struct {
	Mfact  int `json:"mfact"`
	Screen struct {
		Size struct {
			Width uint `json:"width"`
			Heigh uint `json:"height"`
		} `json:"size"`
	} `json:"screen"`

	decor config.Decorations

	master *masterStackWindow
	stack  []masterStackWindow
}

func init() {
	Layouts["master_stack"] = GetMasterStack
}

func GetMasterStack(config []byte, decor config.Decorations) (Layout, error) {
	layout := &MasterStack{
		decor: decor,
	}
	layout.decor.GapsIn /= 2

	if err := json.Unmarshal(config, &layout); err != nil {
		return layout, err
	}

	if layout.Mfact < -100 || layout.Mfact > 100 {
		return layout, errors.New("master stack mfact must be in range -100 to 100")
	} else if layout.Mfact == 0 {
		return layout, errors.New("master stack mfact must not be 0")
	}

	return layout, nil
}

func (l *MasterStack) Map(x *xgbutil.XUtil, e xevent.MapRequestEvent) error {
	window := masterStackWindow{
		id: e.Window,
	}

	if l.master == nil {
		l.master = &window

		if err := xproto.ConfigureWindowChecked(
			x.Conn(), l.master.id,
			xproto.ConfigWindowX|
				xproto.ConfigWindowY|
				xproto.ConfigWindowWidth|
				xproto.ConfigWindowHeight,
			[]uint32{
				uint32(l.decor.GapsOut),
				uint32(l.decor.GapsOut),
				uint32(l.Screen.Size.Width - (2 * l.decor.GapsOut)),
				uint32(l.Screen.Size.Heigh - (2 * l.decor.GapsOut)),
			},
		).Check(); err != nil {
			return fmt.Errorf("Error resizing window: %v\n", err)
		}
	} else if len(l.stack) == 0 {
		l.stack = append(l.stack, window)

		width := ((l.Mfact * int(l.Screen.Size.Width-(2*l.decor.GapsIn)-(2*l.decor.GapsOut))) / 100)

		posX := int(l.Screen.Size.Width) - width - int(l.decor.GapsOut)
		masterPosX := int(l.decor.GapsOut)
		if width < 0 {
			posX = int(l.decor.GapsOut)
			width *= -1
			masterPosX = int(l.Screen.Size.Width) - width - int(l.decor.GapsOut)
		}

		if err := xproto.ConfigureWindowChecked(
			x.Conn(), l.master.id,
			xproto.ConfigWindowX|
				xproto.ConfigWindowY|
				xproto.ConfigWindowWidth|
				xproto.ConfigWindowHeight,
			[]uint32{
				uint32(masterPosX),
				uint32(l.decor.GapsOut),
				uint32((l.Screen.Size.Width - (2 * l.decor.GapsOut) - (2 * l.decor.GapsIn)) - uint(width)),
				uint32(l.Screen.Size.Heigh - (2 * l.decor.GapsOut)),
			},
		).Check(); err != nil {
			return fmt.Errorf("Error resizing window: %v\n", err)
		}

		if err := xproto.ConfigureWindowChecked(
			x.Conn(), window.id,
			xproto.ConfigWindowX|
				xproto.ConfigWindowY|
				xproto.ConfigWindowWidth|
				xproto.ConfigWindowHeight,
			[]uint32{
				uint32(posX),
				uint32(l.decor.GapsOut),
				uint32(width),
				uint32(l.Screen.Size.Heigh - (2 * l.decor.GapsOut)),
			},
		).Check(); err != nil {
			return fmt.Errorf("Error resizing window: %v\n", err)
		}
	} else {
		l.stack = append(l.stack, window)

		height := (l.Screen.Size.Heigh - (2 * l.decor.GapsOut) - uint((len(l.stack)-1)*int(2*l.decor.GapsIn))) / uint(len(l.stack))
		width := ((l.Mfact * int(l.Screen.Size.Width-(2*l.decor.GapsOut)-(2*l.decor.GapsIn))) / 100)

		posX := int(l.Screen.Size.Width) - width - int(l.decor.GapsOut)
		if width < 0 {
			posX = int(l.decor.GapsOut)
			width *= -1
		}

		for i, win := range l.stack {
			if err := xproto.ConfigureWindowChecked(
				x.Conn(), win.id,
				xproto.ConfigWindowX|
					xproto.ConfigWindowY|
					xproto.ConfigWindowWidth|
					xproto.ConfigWindowHeight,
				[]uint32{
					uint32(posX),
					uint32((int(height+(2*l.decor.GapsIn)) * i) + int(l.decor.GapsOut)),
					uint32(width),
					uint32(height),
				},
			).Check(); err != nil {
				return fmt.Errorf("Error resizing window: %v\n", err)
			}
		}
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

			if err := focus(x, e.Event); err != nil {
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
			l.master = &l.stack[0]
		}
		l.stack = []masterStackWindow{}

		if err := xproto.ConfigureWindowChecked(
			x.Conn(), l.master.id,
			xproto.ConfigWindowX|
				xproto.ConfigWindowY|
				xproto.ConfigWindowWidth|
				xproto.ConfigWindowHeight,
			[]uint32{
				uint32(l.decor.GapsOut),
				uint32(l.decor.GapsOut),
				uint32(l.Screen.Size.Width - (2 * l.decor.GapsOut)),
				uint32(l.Screen.Size.Heigh - (2 * l.decor.GapsOut)),
			},
		).Check(); err != nil {
			return fmt.Errorf("Error resizing window: %v\n", err)
		}
	} else if e.Window == l.master.id {
		l.master = &l.stack[0]
		l.stack = l.stack[1:]

		height := (l.Screen.Size.Heigh - (2 * l.decor.GapsOut) - uint((len(l.stack)-1)*int(2*l.decor.GapsIn))) / uint(len(l.stack))
		width := ((l.Mfact * int(l.Screen.Size.Width-(2*l.decor.GapsOut)-(2*l.decor.GapsIn))) / 100)

		posX := int(l.Screen.Size.Width) - width - int(l.decor.GapsOut)
		masterPosX := int(l.decor.GapsOut)
		if width < 0 {
			posX = int(l.decor.GapsOut)
			width *= -1
			masterPosX = int(l.Screen.Size.Width) - width - int(l.decor.GapsOut)
		}

		if err := xproto.ConfigureWindowChecked(
			x.Conn(), l.master.id,
			xproto.ConfigWindowX|
				xproto.ConfigWindowY|
				xproto.ConfigWindowWidth|
				xproto.ConfigWindowHeight,
			[]uint32{
				uint32(masterPosX),
				uint32(l.decor.GapsOut),
				uint32(l.Screen.Size.Width - uint(width) - (2 * l.decor.GapsIn) - (2 * l.decor.GapsOut)),
				uint32(l.Screen.Size.Heigh - (2 * l.decor.GapsOut)),
			},
		).Check(); err != nil {
			return fmt.Errorf("Error resizing window: %v\n", err)
		}

		for i, win := range l.stack {
			if err := xproto.ConfigureWindowChecked(
				x.Conn(), win.id,
				xproto.ConfigWindowX|
					xproto.ConfigWindowY|
					xproto.ConfigWindowWidth|
					xproto.ConfigWindowHeight,
				[]uint32{
					uint32(posX),
					uint32((int(height+(2*l.decor.GapsIn)) * i) + int(l.decor.GapsOut)),
					uint32(width),
					uint32(height),
				},
			).Check(); err != nil {
				return fmt.Errorf("Error resizing window: %v\n", err)
			}
		}
	} else if len(l.stack) == 1 {
		l.stack = []masterStackWindow{}

		if err := xproto.ConfigureWindowChecked(
			x.Conn(), l.master.id,
			xproto.ConfigWindowX|
				xproto.ConfigWindowY|
				xproto.ConfigWindowWidth|
				xproto.ConfigWindowHeight,
			[]uint32{
				uint32(0),
				uint32(0),
				uint32(l.Screen.Size.Width),
				uint32(l.Screen.Size.Heigh),
			},
		).Check(); err != nil {
			return fmt.Errorf("Error resizing window: %v\n", err)
		}
	} else {
		for i, win := range l.stack {
			if win.id == e.Window {
				l.stack = append(l.stack[:i], l.stack[i+1:]...)
			}
		}

		height := (l.Screen.Size.Heigh - (2 * l.decor.GapsOut) - uint((len(l.stack)-1)*int(2*l.decor.GapsIn))) / uint(len(l.stack))
		width := ((l.Mfact * int(l.Screen.Size.Width-(2*l.decor.GapsOut)-(2*l.decor.GapsIn))) / 100)

		posX := int(l.Screen.Size.Width) - width - int(l.decor.GapsOut)
		if width < 0 {
			posX = int(l.decor.GapsOut)
			width *= -1
		}

		for i, win := range l.stack {
			if err := xproto.ConfigureWindowChecked(
				x.Conn(), win.id,
				xproto.ConfigWindowX|
					xproto.ConfigWindowY|
					xproto.ConfigWindowWidth|
					xproto.ConfigWindowHeight,
				[]uint32{
					uint32(posX),
					uint32((int(height+(2*l.decor.GapsIn)) * i) + int(l.decor.GapsOut)),
					uint32(width),
					uint32(height),
				},
			).Check(); err != nil {
				return fmt.Errorf("Error resizing window: %v\n", err)
			}
		}
	}

	xevent.Detach(x, e.Window)
	keybind.Detach(x, e.Window)
	mousebind.Detach(x, e.Window)

	return nil
}

func focus(x *xgbutil.XUtil, id xproto.Window) error {
	return xproto.SetInputFocusChecked(
		x.Conn(),
		xproto.InputFocusPointerRoot,
		id,
		xproto.TimeCurrentTime,
	).Check()
}
