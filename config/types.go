package config

type Config struct {
	Keybinds []Keybind `json:"keybinds"`
	Layout   Layout    `json:"layout"`
}

type Keybind struct {
	Cmd  []string `json:"cmd"`
	Bind string   `json:"bind"`
}

type Layout struct {
	Name   string         `json:"name"`
	Config map[string]any `json:"config"`
}
