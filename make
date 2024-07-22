#!/bin/sh

CURRENT_DIR=`pwd`
SCRIPT_DIR=`dirname "$0"`

cd $SCRIPT_DIR
go run ./cmd/make $@
cd $CURRENT_DIR
