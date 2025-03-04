package main

import (
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplication_setupFunctions(t *testing.T) {
	a := Application{}
	h := cmd.Handlers{
		{
			Func: "portald",
			Name: "portald",
		},
		{
			Func: "app",
			Name: "app",
		},
	}
	a.setupFunctions(h)
	assert.True(t, "portald" != h[0].Func)
	assert.True(t, "app" != h[1].Func)
	assert.True(t, h[0].Func != h[1].Func)
}
