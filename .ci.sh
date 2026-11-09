#! /bin/sh

set -eu

run() {
  echo "=> ${@}"

  if output=$("${@}" 2>&1); then
    return 0
  else
    status=${?}
    echo "${output}"
    return "${status}"
  fi
}

if test -n "${CI:-}"; then
  run gofmt -d .
else
  run go mod tidy
  run go fmt ./...
fi

run go vet "${@}" ./...
run go fix -diff "${@}" ./...
run go test "${@}" ./...
