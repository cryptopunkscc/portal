package debug

import (
	"testing"
	"time"

	api "github.com/cryptopunkscc/portal/cmd/astral-yt-dlp/api"
	"github.com/cryptopunkscc/portal/cmd/astral-yt-dlp/src"
	"github.com/cryptopunkscc/portal/pkg/test"
)

func (c *TestContext) TestServeAstralYouTubeDl(requests ...api.Request) test.Test {
	return c.Test().Args(requests).Func(func(t *testing.T) {
		p := api.Client{Client: c.Apphost.Client}
		for _, request := range requests {
			err := p.Download(c.Context, request)
			test.NoError(t, err)
		}
		ch, erp := p.Status(c.Context)
		test.NoError(t, erp)
		for progress := range ch {
			t.Log(*progress)
		}
		test.NoError(t, erp)
	}).Requires(
		c.ServeAstralYouTubeDlService(),
	)
}

func (c *TestContext) ServeAstralYouTubeDlService() test.Test {
	return c.Test().Func(func(t *testing.T) {
		go func() {
			s := astral_yt_dlp.Service{Dir: "./local"}
			err := s.Serve(c.Context)
			test.NoError(t, err)
		}()
		time.Sleep(time.Second)
	}).Requires(
		c.CreateUser(),
	)
}
