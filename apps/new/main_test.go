package main

import (
	gotest "github.com/cryptopunkscc/portal/target/go/test"
	"testing"
)

func TestResolve(t *testing.T) {
	gotest.ResolveGoProject(t)
}

func TestBuild(t *testing.T) {
	gotest.BuildGoProject(t)
}
