package main

import (
	"context"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"testing"
	"time"
)

func init() {
	plog.Verbosity = 100
}

func TestApplication_start(t *testing.T) {
	dir := test.Mkdir(t)
	c := portal.Config{}
	c.Dir = dir
	c.Node.Log.Level = 100
	c.Apphost.Listen = []string{"tcp:127.0.0.1:8635"}
	if err := writeConfig(t, c, dir, portal.DefaultConfigFile); err != nil {
		plog.P().Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		time.Sleep(100 * time.Millisecond) // give a time to kill astrald process
	})
	args := RunArgs{ConfigPath: dir}
	err := testApplication().start(ctx, args)
	test.AssertErr(t, err)
}

func TestApplication_start_project_config(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		time.Sleep(100 * time.Millisecond) // give a time to kill astrald process
	})
	args := RunArgs{}
	err := testApplication().start(ctx, args)
	test.AssertErr(t, err)
}

func testApplication() (a *Application[target.Portal_]) {
	a = &Application[target.Portal_]{}
	a.ExtraTokens = []string{"portal"}
	return
}
