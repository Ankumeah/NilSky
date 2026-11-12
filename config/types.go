package config

type Config struct {
	Keybinds    []Keybind   `json:"keybinds"`
	Layout      Layout      `json:"layout"`
	Decorations Decorations `json:"decorations"`
}

type Keybind struct {
	Cmd  []string `json:"cmd"`
	Bind string   `json:"bind"`
}

type Layout struct {
	Name   string         `json:"name"`
	Config map[string]any `json:"config"`
}

type Decorations struct {
	GapsIn      uint   `json:"gaps_in"`
	GapsOut     uint   `json:"gaps_out"`
	BorderWidth uint   `json:"border_width"`
	BorderColor string `json:"border_color"`
}
