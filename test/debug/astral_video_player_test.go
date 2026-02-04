package debug

import (
	"path"
	"testing"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/fs"
	"github.com/cryptopunkscc/portal/apps/player/video"
	"github.com/cryptopunkscc/portal/apps/player/vlc"
	"github.com/cryptopunkscc/portal/cmd/astral-yt-dlp/api"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/require"
)

func (c *TestContext) ServeAstralVideoPlayer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		go func() {
			var err error
			s := video.Service{}
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
		time.Sleep(time.Second)
		audioFileId := c.findVideoFile(t)
		c.TestMediaPlayerClient(t, "video", audioFileId)
	}).Requires(
		c.NewWatch(),
		c.ServeAstralVideoPlayer(),
		c.TestServeAstralYouTubeDl(
			astral_yt_dlp.Request{Url: "https://www.youtube.com/watch?v=Kt7ZDFKFNxc", Audio: false},
		),
	)
}

func (c *TestContext) findVideoFile(t *testing.T) (audioFileId *astral.ObjectID) {
	scan, errPtr := c.Apphost.Objects().ObjectsClient.Scan(c.Context, "test", false)
	for id := range scan {
		descCh, errPtr := c.Apphost.Objects().ObjectsClient.Describe(c.Context, id)
		for desc := range descCh {
			if fl, ok := desc.Descriptor.(*fs.FileLocation); ok {
				if path.Ext(fl.Path.String()) == ".mkv" {
					audioFileId = id
				}
			}
		}
		test.NoError(t, *errPtr)
	}
	test.NoError(t, *errPtr)
	require.NotNil(t, audioFileId)
	return
}
