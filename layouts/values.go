package layouts

var Layouts = map[string]func(config []byte) (Layout, error){}
