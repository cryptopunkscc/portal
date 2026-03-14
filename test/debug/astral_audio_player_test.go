package debug

import (
	"testing"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/media"
	"github.com/cryptopunkscc/portal/cmd/astral-yt-dlp/client"
	"github.com/cryptopunkscc/portal/pkg/util/player/audio"
	"github.com/cryptopunkscc/portal/pkg/util/player/beep"
	"github.com/cryptopunkscc/portal/pkg/util/test"
	"github.com/stretchr/testify/require"
)

func (c *TestContext) ServeAstralAudioPlayer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		go func() {
			s := audio.Service{}
			s.Player = &beep.Player{}
			err := s.Serve(c.Context)
			test.NoError(t, err)
		}()
		time.Sleep(time.Second)
	})
}

func (c *TestContext) TestAstralAudioPlayer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		audioFileId := c.findAudioFile(t)
		c.TestMediaPlayerClient(t, "audio", audioFileId)
	}).Requires(
		c.NewWatch(),
		c.ServeAstralAudioPlayer(),
		c.TestServeAstralYouTubeDl(
			astral_yt_dlp.Request{Url: "https://www.youtube.com/watch?v=YAszKsWBpKs", Audio: true},
		),
	)
}

func (c *TestContext) findAudioFile(t *testing.T) (audioFileId *astral.ObjectID) {
	scan, errPtr := c.Client.Objects().Scan(c.Context, "test", false)
	for id := range scan {
		descCh, errPtr := c.Client.Objects().Describe(c.Context, id)
		for desc := range descCh {
			if _, ok := desc.Descriptor.(*media.AudioFile); ok {
				audioFileId = id
			}
		}
		test.NoError(t, errPtr)
	}
	test.NoError(t, errPtr)
	require.NotNil(t, audioFileId)
	return
}
