package layouts

import (
	"fmt"

	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/ewmh"
	"github.com/jezek/xgbutil/xwindow"
)

func (*MasterStack) focus(x *xgbutil.XUtil, id xproto.Window) error {
	return xproto.SetInputFocusChecked(
		x.Conn(),
		xproto.InputFocusPointerRoot,
		id,
		xproto.TimeCurrentTime,
	).Check()
}

func (l *MasterStack) redraw(x *xgbutil.XUtil) error {
	ids := []xproto.Window{}
	for _, win := range append(l.stack, l.master) {
		if win != nil {
			ids = append(ids, win.id)
		}
	}
	ewmh.ClientListSet(x, ids)

	if l.master == nil {
		return nil
	}

	if err := l.newFramedWindow(x, l.master); err != nil {
		return err
	}

	for _, win := range l.stack {
		if err := l.newFramedWindow(x, win); err != nil {
			return err
		}
	}

	return nil
}

func (l *MasterStack) newFramedWindow(x *xgbutil.XUtil, win *masterStackWindow) error {
	if win.frameID == 0 {
		frame, err := xwindow.Generate(x)
		if err != nil {
			return err
		}

		if err := frame.CreateChecked(
			x.RootWin(),
			int(win.x),
			int(win.y),
			int(win.width),
			int(win.height),
			xproto.CwBackPixel,
			uint32(l.borderColor),
		); err != nil {
			return err
		}

		if err := frame.Listen(xproto.EventMaskSubstructureNotify); err != nil {
			return fmt.Errorf("Error listing for EventMaskSubstructureNotify: %w\n", err)
		}

		if err := xproto.ChangeSaveSetChecked(
			x.Conn(),
			xproto.SetModeInsert,
			win.id,
		).Check(); err != nil {
			return err
		}

		if err := xproto.ReparentWindowChecked(
			x.Conn(),
			win.id,
			frame.Id,
			int16(l.decor.BorderWidth),
			int16(l.decor.BorderWidth),
		).Check(); err != nil {
			return err
		}

		frame.Map()

		win.frameID = frame.Id
	} else {
		if err := xproto.ConfigureWindowChecked(
			x.Conn(), win.frameID,
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

	if err := xproto.ConfigureWindowChecked(
		x.Conn(), win.id,
		xproto.ConfigWindowX|
			xproto.ConfigWindowY|
			xproto.ConfigWindowWidth|
			xproto.ConfigWindowHeight,
		[]uint32{
			uint32(l.decor.BorderWidth),
			uint32(l.decor.BorderWidth),
			uint32(win.width - (2 * l.decor.BorderWidth)),
			uint32(win.height - (2 * l.decor.BorderWidth)),
		},
	).Check(); err != nil {
		return fmt.Errorf("Error resizing window: %v\n", err)
	}

	return nil
}
