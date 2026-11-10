package main

import (
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
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
		[]uint32{
			xproto.EventMaskSubstructureRedirect |
				xproto.EventMaskSubstructureNotify,
		},
	).Check(); err != nil {
		return nil, err
	}

	return x, nil
}
