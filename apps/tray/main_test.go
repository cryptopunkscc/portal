package main

import (
	"testing"

	gotest "github.com/cryptopunkscc/portal/target/go/test"
)

func TestResolve(t *testing.T) {
	gotest.ResolveGoProject(t)
}

func TestBuild(t *testing.T) {
	gotest.BuildGoProject(t)
}
