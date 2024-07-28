package test

import (
	"github.com/cryptopunkscc/portal/target2"
	golang "github.com/cryptopunkscc/portal/target2/go"
	"testing"
)

func TestList(t *testing.T) {
	target2.Any[target2.Base](
		target2.Try(golang.ResolveProject),
	)
}
