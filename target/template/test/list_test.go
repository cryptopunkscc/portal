package test

import (
	"github.com/cryptopunkscc/portal/target/template"
	"testing"
)

func TestList(t *testing.T) {
	list := template.List()
	str := list.MarshalCLI()
	println(str)
}
