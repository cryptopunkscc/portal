#!/bin/bash

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

go run \
  -tags dev \
  -gcflags "all=-N -l" \
  "$script_dir" "$@"