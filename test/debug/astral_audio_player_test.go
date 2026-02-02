package debug

import (
	"testing"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/media"
	astral_audio_player2 "github.com/cryptopunkscc/portal/cmd/astral-audio-player/client"
	"github.com/cryptopunkscc/portal/cmd/astral-audio-player/src"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/require"
)

func (c *TestContext) ServeAstralAudioPlayer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		go func() {
			err := astral_audio_player.Serve(c.Context)
			test.NoError(t, err)
		}()
		time.Sleep(time.Second)
	})
}

func (c *TestContext) TestAstralAudioPlayer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		debugDelay := time.Second * 0
		// find file to play
		audioFileId := c.findAudioFile(t)

		// play in the background
		audioClient := astral_audio_player2.Client{Client: c.Apphost.Client.WithTarget(c.Apphost.HostID())}
		err := audioClient.PlayID(c.Context, *audioFileId)
		test.NoError(t, err)
		time.Sleep(debugDelay)

		time.Sleep(time.Millisecond * 100)

		err = audioClient.Pause(c.Context)
		test.NoError(t, err)
		time.Sleep(time.Second)

		// move 30s forward
		err = audioClient.Add(c.Context, 60*time.Second)
		test.NoError(t, err)

		err = audioClient.Resume(c.Context)
		test.NoError(t, err)

		status, err := audioClient.Status(c.Context)
		test.NoError(t, err)
		t.Logf("status: %+v", status)

		// move back to start
		time.Sleep(debugDelay)
		time.Sleep(time.Millisecond)
		err = audioClient.Seek(c.Context, 0)

		// close
		time.Sleep(debugDelay)
		time.Sleep(time.Millisecond)
		err = audioClient.Stop(c.Context)

	}).Requires(
		c.NewWatch(),
		c.ServeAstralAudioPlayer(),
	)
}

func (c *TestContext) findAudioFile(t *testing.T) (audioFileId *astral.ObjectID) {
	scan, errPtr := c.Apphost.Objects().ObjectsClient.Scan(c.Context, "test", false)
	for id := range scan {
		descCh, errPtr := c.Apphost.Objects().ObjectsClient.Describe(c.Context, id)
		for desc := range descCh {
			if _, ok := desc.Descriptor.(*media.AudioFile); ok {
				audioFileId = id
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
