package layouts

import (
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/xevent"
)

type Layout interface {
	AddWindow(x *xgbutil.XUtil, e xevent.MapRequestEvent) error
	RemoveWindow(x *xgbutil.XUtil, e xevent.MapRequestEvent) error
}
