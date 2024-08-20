package factory

import _ "github.com/cryptopunkscc/astrald/mod/apphost/src"
import (
	"github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runtime/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuntime(t *testing.T) {
	api := testApi{T: t, events: make(chan *mobile.Event, 256)}
	portal := client.Portal(target.PortPortal.String())
	runtime := Runtime(api)
	runtime.Start()
	assert.Equal(t, &mobile.Event{Msg: mobile.STARTING}, <-api.events)
	t.Log("STARTING")
	assert.Equal(t, &mobile.Event{Msg: mobile.STARTED}, <-api.events)
	t.Log("STARTED")
	if err := portal.Ping(); err != nil {
		t.Fatal(err)
	}
	runtime.Stop()
	assert.Equal(t, &mobile.Event{Msg: mobile.STOPPED}, <-api.events)
	t.Log("STOPPED")
}

type testApi struct {
	*testing.T
	events chan *mobile.Event
}

func (t testApi) NodeRoot() string          { return "test_node_root" }
func (t testApi) RequestHtml(string) error  { return nil }
func (t testApi) Event(event *mobile.Event) { t.events <- event }
