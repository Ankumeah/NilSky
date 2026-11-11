package layouts

import (
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/xevent"
)

type Layout interface {
	Map(x *xgbutil.XUtil, e xevent.MapRequestEvent) error
}
