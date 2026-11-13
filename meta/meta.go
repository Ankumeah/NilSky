package meta

import "path/filepath"

const LowerName = "nilsky"

var SocketFile = LowerName + ".sock"
var ConfigFile = filepath.Join(LowerName, "config.json")
