# NilSky

NilSky is a minimal [x11](https://en.wikipedia.org/wiki/X_Window_System) WM written in [Golang](https://go.dev)

## Config
NilSky `$XDG_CONFIG_HOME/nilsky/config.json` or `$HOME/.config/nilsky/config.json` in order

A spec is provided [here](./config_spec.yaml)
> Side note: I had no idea we could do YAML written specs. Hell yeah

An example config is provided [here](./example.config.json)
which lists every single config option

## IPC
NilSky listen on a [unix socket](https://en.wikipedia.org/wiki/Unix_domain_socket)
located at `$XDG_RUNTIME_DIR/nilsky.sock` to provide easy direct communication with the WM

Given below are the commands supported by the socket

> [!NOTE]
> All successful commands start with `OK: ` and all errors start with `ERROR: `

### `ping`
Just responds with `OK: PONG`

#### Examples
- `echo 'ping' | nc -q 0 -U /run/user/1000/nilsky.sock`

### `close_window <window_pid | focus>`
Send a close request to the window specified by window_pid or the currently focused window if focus is the arg

#### Examples
- `echo 'close_window focus' | nc -q 0 -U /run/user/1000/nilsky.sock`
- `echo 'close_window 42042' | nc -q 0 -U /run/user/1000/nilsky.sock`

## Lastest release
[v0.1.0](https://github.com/Ankumeah/NilSky/releases/tag/v0.1.0)

## Fun side note
For dev, I kinda had to use a cursed looking config to make sure
I can properly see the border and gaps and to make sure [kitty](https://sw.kovidgoyal.net/kitty/)
lunches within [Xephyr](https://freedesktop.org/wiki/Software/Xephyr) and not my host WM

If you want, you can try [it](./dev.config.json) out

## Depends
- `github.com/jezek/xgb` and `github.com/jezek/xgbutil`: by [jezek](https://github.com/jezek) (Thank you for this wonderful lib)
