package config

type Config struct {
	Keybinds []Keybind `json:"keybinds"`
}

type Keybind struct {
	Cmd  []string `json:"cmd"`
	Bind string   `json:"bind"`
}
