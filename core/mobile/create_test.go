package core

import (
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/test"
	"log"
	"testing"
)

func TestCreate(t *testing.T) {
	dir := test.Dir(t)
	env.AstraldHome.SetDir(dir)
	env.AstraldDb.SetDir(dir)
	env.PortaldApps.SetDir(dir, "apps")
	env.PortaldTokens.SetDir(dir, "tokens")

	core := Create(testApi{})
	s, ok := core.(*service)
	if !ok {
		t.Fatal()
	}
	log.Println(s)
}

type testApi struct{ testDir string }

func (t testApi) CacheDir() string { return t.testDir }

func (t testApi) DataDir() string { return t.testDir }

func (t testApi) DbDir() string {
	//TODO implement me
	return t.testDir
}

func (t testApi) Event(event *mobile.Event) {
	//TODO implement me
	panic("implement me")
}

func (t testApi) StartHtml(pkg string, args string) error {
	//TODO implement me
	panic("implement me")
}

func (t testApi) Net() mobile.Net {
	//TODO implement me
	panic("implement me")
}

var _ mobile.Api = testApi{}
