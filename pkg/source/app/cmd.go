package app

import (
	"context"
	"io/fs"

	"github.com/cryptopunkscc/portal/pkg/client"
	"github.com/cryptopunkscc/portal/pkg/source"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
)

type Publisher struct {
	*client.Objects
}

// PublishBundles recursively searches the given path and publishes any app bundle it finds.
func (p Publisher) PublishBundles(ctx context.Context, path string) (out []ReleaseInfo, err error) {
	return p.PublishBundlesSrc(ctx, source.OSRef(path))
}

func (p Publisher) PublishBundlesSrc(ctx context.Context, src source.Source) (out []ReleaseInfo, err error) {
	if p.Objects == nil {
		p.Objects = client.Default.Objects()
	}
	defer plog.TraceErr(&err)
	apps := source.CollectIt(src, &Bundle{})
	if len(apps) == 0 {
		panic(src.Ref_().Path)
		return nil, fs.ErrNotExist
	}

	for _, app := range apps {
		var info ReleaseInfo
		if info, err = app.Publish(ctx, p.Objects); err != nil {
			return
		}
		out = append(out, info)
	}
	return
}
