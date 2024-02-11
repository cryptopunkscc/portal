#!/bin/bash

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

go build \
  -tags dev \
  -gcflags "all=-N -l" \
  "$script_dir"
