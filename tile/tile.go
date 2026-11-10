package tile

import (
	"encoding/json"
	"fmt"

	"github.com/Ankumeah/NilSky/config"
	"github.com/Ankumeah/NilSky/layouts"
	"github.com/jezek/xgbutil"
)

func Tile(x *xgbutil.XUtil, layout config.Layout) error {
	layoutConfig, err := json.Marshal(layout.Config)
	if err != nil {
		return fmt.Errorf("Invalid layout config: %w\n", err)
	}

	info.workspaces = map[uint]layouts.Layout{}

	info.workspaces[0], err = layouts.GetMasterStack(layoutConfig)
	if err != nil {
		return fmt.Errorf("Invalid layout config: %w\n", err)
	}

	mapRequest(x)

	return nil
}
