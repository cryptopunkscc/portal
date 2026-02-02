package debug

import (
	"testing"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	client "github.com/cryptopunkscc/portal/apps/player/client"
	"github.com/cryptopunkscc/portal/pkg/test"
)

func (c *TestContext) TestMediaPlayerClient(
	t *testing.T,
	service string,
	fileId *astral.ObjectID,
) {
	debugDelay := time.Second * 2
	p := client.Client{
		Name:   service,
		Client: c.Apphost.Client.WithTarget(c.Apphost.HostID()),
	}

	t.Run("play", func(t *testing.T) {
		err := p.PlayID(c.Context, *fileId)
		test.NoError(t, err)
	})

	time.Sleep(debugDelay)
	time.Sleep(time.Millisecond * 100)

	t.Run("pause", func(t *testing.T) {
		err := p.Pause(c.Context)
		test.NoError(t, err)
	})

	time.Sleep(time.Second)

	t.Run("move", func(t *testing.T) {
		err := p.Move(c.Context, 60*time.Second)
		test.NoError(t, err)
	})

	t.Run("status", func(t *testing.T) {
		status, err := p.Status(c.Context)
		test.NoError(t, err)
		t.Logf("status: %+v", status)
	})

	t.Run("resume", func(t *testing.T) {
		err := p.Resume(c.Context)
		test.NoError(t, err)
	})

	time.Sleep(debugDelay)
	time.Sleep(time.Millisecond)

	t.Run("seek", func(t *testing.T) {
		err := p.Seek(c.Context, 0)
		test.NoError(t, err)
	})

	time.Sleep(debugDelay)
	time.Sleep(time.Millisecond)

	t.Run("stop", func(t *testing.T) {
		err := p.Stop(c.Context)
		test.NoError(t, err)
	})
}
