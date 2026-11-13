package tile

import (
	"encoding/json"
	"fmt"

	"github.com/Ankumeah/NilSky/config"
	"github.com/Ankumeah/NilSky/layouts"
	"github.com/jezek/xgbutil"
)

func Tile(x *xgbutil.XUtil, config config.Config) error {
	layoutConfig, err := json.Marshal(config.Layout.Config)
	if err != nil {
		return fmt.Errorf("Invalid layout config: %w\n", err)
	}

	info.workspaces = map[uint]layouts.Layout{}

	layout, ok := layouts.Layouts[config.Layout.Name]
	if !ok {
		return fmt.Errorf("Unknown layout: %w\n", err)
	}

	info.workspaces[0], err = layout(layoutConfig, config.Decorations)
	if err != nil {
		return fmt.Errorf("Invalid layout config: %w\n", err)
	}

	mapRequest(x)

	return nil
}
