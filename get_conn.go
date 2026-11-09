package main

import (
	"log"

	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/xevent"
)

func GetConn() (*xgbutil.XUtil, error) {
	x, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}

	if err = xproto.ChangeWindowAttributesChecked(
		x.Conn(),
		x.RootWin(),
		xproto.CwEventMask,
		[]uint32{xproto.EventMaskSubstructureRedirect | xproto.EventMaskSubstructureNotify},
	).Check(); err != nil {
		return nil, err
	}

	xevent.MapRequestFun(
		func(x *xgbutil.XUtil, e xevent.MapRequestEvent) {
			xproto.MapWindow(x.Conn(), e.Window)

			screen, err := xproto.GetGeometry(
				x.Conn(),
				xproto.Drawable(x.RootWin()),
			).Reply()
			if err != nil {
				log.Printf("Error getting screen size: %v\n", err)
				return
			}

			if err = xproto.ConfigureWindowChecked(
				x.Conn(), e.Window,
				xproto.ConfigWindowWidth|xproto.ConfigWindowHeight,
				[]uint32{uint32(screen.Width), uint32(screen.Height)},
			).Check(); err != nil {
				log.Printf("Error resizing window: %v\n", err)
				return
			}

			if err := xproto.MapWindowChecked(x.Conn(), e.Window).Check(); err != nil {
				log.Printf("Error mapping window: %v\n", err)
				return
			}
		},
	).Connect(x, x.RootWin())

	return x, nil
}
