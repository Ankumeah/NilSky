package tile

import (
	"log"

	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/xevent"
)

func mapRequest(x *xgbutil.XUtil) {
	xevent.MapRequestFun(
		func(x *xgbutil.XUtil, e xevent.MapRequestEvent) {
			info.mu.Lock()
			defer info.mu.Unlock()

			if err := info.workspaces[info.currentWS].AddWindow(x, e); err != nil {
				log.Printf("Error mapping window: %v\n", err)
			}
		},
	).Connect(x, x.RootWin())
}
