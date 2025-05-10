package basic

import (
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/js"
	"github.com/cryptopunkscc/portal/target/source"
	"testing"
)

func TestResolve(t *testing.T) {
	s := source.Dir("js")
	_, err := js.ResolveDist(s)
	test.AssertErr(t, err)
}
