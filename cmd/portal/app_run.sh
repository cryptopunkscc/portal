#!/bin/bash

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

go run \
  -tags "desktop,wv2runtime.download,production" \
  -ldflags "-w -s" \
  "$script_dir" "$@"
