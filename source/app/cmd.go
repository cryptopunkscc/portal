package app

import (
	"io/fs"

	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
)

// PublishAppBundles recursively searches the given path and publishes any app bundle it finds.
func PublishAppBundles(path string) (out []ReleaseInfo, err error) {
	return PublishAppBundlesSrc(astrald.DefaultClient(), source.OSRef(path))
}

func PublishAppBundlesSrc(client *astrald.Client, src source.Source) (out []ReleaseInfo, err error) {
	defer plog.TraceErr(&err)
	apps := source.CollectIt(src, &Bundle{})
	if len(apps) == 0 {
		return nil, fs.ErrNotExist
	}
	objects := astrald.NewObjectsClient(client, "")
	for _, app := range apps {
		var info ReleaseInfo
		if info, err = app.Publish(objects); err != nil {
			return
		}
		out = append(out, info)
	}
	return
}
