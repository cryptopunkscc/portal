package template

import (
	"testing"
)

func TestList(t *testing.T) {
	list := List()
	str := list.MarshalCLI()
	println(str)
}
