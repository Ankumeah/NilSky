// [WARNING]
//
// THIS IS PROBABLY ONE OF THE
// UGLIEST FIELS I HAVE EVER WRITTEN
// AND IV WROTE A LOT OF UGLY FIELS
//
// YOU HAVE BEEN WARNED

package layouts

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/xevent"
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

	master *masterStackWindow
	stack  []masterStackWindow
}

func init() {
	Layouts["master_stack"] = GetMasterStack
}

func GetMasterStack(config []byte) (Layout, error) {
	layout := &MasterStack{}

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

func (l *MasterStack) AddWindow(x *xgbutil.XUtil, e xevent.MapRequestEvent) error {
	window := masterStackWindow{
		id: e.Window,
	}

	if l.master == nil {
		l.master = &window

		if err := xproto.ConfigureWindowChecked(
			x.Conn(), l.master.id,
			xproto.ConfigWindowWidth|xproto.ConfigWindowHeight,
			[]uint32{
				uint32(l.Screen.Size.Width),
				uint32(l.Screen.Size.Heigh),
			},
		).Check(); err != nil {
			return fmt.Errorf("Error resizing window: %v\n", err)
		}
	} else if len(l.stack) == 0 {
		l.stack = append(l.stack, window)

		width := (l.Mfact * int(l.Screen.Size.Width)) / 100

		posX := int(l.Screen.Size.Width) - width
		masterPosX := 0
		if width < 0 {
			posX = 0
			width *= -1
			masterPosX = int(l.Screen.Size.Width) - width
		}

		if err := xproto.ConfigureWindowChecked(
			x.Conn(), l.master.id,
			xproto.ConfigWindowX|
				xproto.ConfigWindowY|
				xproto.ConfigWindowWidth|
				xproto.ConfigWindowHeight,
			[]uint32{
				uint32(masterPosX),
				uint32(0),
				uint32(l.Screen.Size.Width - uint(width)),
				uint32(l.Screen.Size.Heigh),
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
				uint32(0),
				uint32(width),
				uint32(l.Screen.Size.Heigh),
			},
		).Check(); err != nil {
			return fmt.Errorf("Error resizing window: %v\n", err)
		}
	} else {
		l.stack = append(l.stack, window)

		height := l.Screen.Size.Heigh / uint(len(l.stack))
		width := (l.Mfact * int(l.Screen.Size.Width)) / 100

		posX := int(l.Screen.Size.Width) - width
		if width < 0 {
			posX = 0
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
					uint32(int(height) * i),
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

	return nil
}

func (l *MasterStack) RemoveWindow(x *xgbutil.XUtil, e xevent.MapRequestEvent) error {
	return errors.New("Unimplimented") // TODO: Do something with my life, idk
}
