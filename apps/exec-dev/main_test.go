package main

import (
	gotest "github.com/cryptopunkscc/portal/target/go/test"
	"testing"
)

func TestResolve(t *testing.T) {
	gotest.ResolveGoProject(t)
}
