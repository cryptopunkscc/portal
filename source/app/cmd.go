package app

import (
	"io/fs"

	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
)

type Publisher struct {
	*astrald.ObjectsClient
}

// PublishBundles recursively searches the given path and publishes any app bundle it finds.
func (p Publisher) PublishBundles(path string) (out []ReleaseInfo, err error) {
	return p.PublishBundlesSrc(source.OSRef(path))
}

func (p Publisher) PublishBundlesSrc(src source.Source) (out []ReleaseInfo, err error) {
	if p.ObjectsClient == nil {
		p.ObjectsClient = astrald.Objects()
	}
	defer plog.TraceErr(&err)
	apps := source.CollectIt(src, &Bundle{})
	if len(apps) == 0 {
		panic(src.Ref_().Path)
		return nil, fs.ErrNotExist
	}

	for _, app := range apps {
		var info ReleaseInfo
		if info, err = app.Publish(p.ObjectsClient); err != nil {
			return
		}
		out = append(out, info)
	}
	return
}
