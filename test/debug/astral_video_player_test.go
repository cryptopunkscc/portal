package debug

import (
	"path"
	"testing"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/fs"
	player "github.com/cryptopunkscc/portal/apps/player/src"
	"github.com/cryptopunkscc/portal/apps/player/vlc"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/require"
)

func (c *TestContext) ServeAstralVideoPlayer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		go func() {
			var err error
			s := player.Service{Name: "video"}
			s.Player, err = vlc.NewPlayer()
			test.NoError(t, err)
			err = s.Serve(c.Context)
			test.NoError(t, err)
		}()
		time.Sleep(time.Second)
	})
}

func (c *TestContext) TestAstralVideoPlayer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		audioFileId := c.findVideoFile(t)
		c.TestMediaPlayerClient(t, "video", audioFileId)
	}).Requires(
		c.NewWatch(),
		c.ServeAstralVideoPlayer(),
	)
}

func (c *TestContext) findVideoFile(t *testing.T) (audioFileId *astral.ObjectID) {
	scan, errPtr := c.Apphost.Objects().ObjectsClient.Scan(c.Context, "test", false)
	for id := range scan {
		descCh, errPtr := c.Apphost.Objects().ObjectsClient.Describe(c.Context, id)
		for desc := range descCh {
			if fl, ok := desc.Descriptor.(*fs.FileLocation); ok {
				if path.Ext(fl.Path.String()) == ".rmvb" {
					audioFileId = id
				}
			}
		}
		if errPtr != nil {
			test.NoError(t, *errPtr)
		}
	}
	if errPtr != nil {
		test.NoError(t, *errPtr)
	}
	require.NotNil(t, audioFileId)
	return
}
