#!/bin/bash

call_dir=`pwd`
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd $script_dir
wails build --debug
cd $call_dir