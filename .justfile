default: build run

build:
  CGO_ENABLED=0 go build -ldflags='-s -w' -o build/out .

run display = ":1" *args:
  #! /bin/sh

  if ! test -e /tmp/.X1-lock; then
    Xephyr :1 -ac -screen 1920x1080 &
    sleep 2
  fi

  DISPLAY={{ display }} ./build/out {{ args }}

ci ci = "":
  #! /bin/sh

  if test -n '{{ ci }}'; then
    echo 'CI = true'
    export CI=true
  fi

  ./.ci.sh
