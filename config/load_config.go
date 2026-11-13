package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Ankumeah/NilSky/meta"
)

func LoadConfig(path string) (Config, error) {
	config := Config{}

	if path == "" {
		var err error
		path, err = os.UserConfigDir()
		if err != nil {
			return config, err
		}

		path = filepath.Join(path, meta.ConfigFile)
	}

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}

	return config, json.NewDecoder(file).Decode(&config)
}
