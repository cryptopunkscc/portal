package apps

import (
	"context"
	"encoding/json"
	"github.com/cryptopunkscc/astrald/object"
	"github.com/cryptopunkscc/portal/api/apps"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
)

type Apps struct {
	Dir     target.Source
	File    target.File
	Find    target.Find[target.Bundle_]
	Resolve target.Resolve[target.App_]
}

func (a Apps) Get(ctx context.Context, pkg string) (out target.App_, err error) {
	list, err := a.List(ctx)
	if err != nil {
		return
	}
	for _, app := range list {
		if filepath.Base(app.Abs()) == pkg {
			return app, nil
		}
		if app.Manifest().Match(pkg) {
			return app, nil
		}
	}
	return nil, target.ErrNotFound
}

func (a Apps) List(_ context.Context) (target.Portals[target.App_], error) {
	return a.Resolve.List(a.Dir), nil
}

func (a Apps) Observe(ctx context.Context) (out <-chan apps.App, err error) {
	list, err := a.list(ctx)
	if err != nil {
		return
	}
	changes, err := a.changes(ctx)
	if err != nil {
		return
	}
	out = flow.Combine[apps.App](list, changes)
	return
}

func (a Apps) list(ctx context.Context) (out <-chan apps.App, err error) {
	list, err := a.List(ctx)
	if err != nil {
		return
	}
	return flow.Map(flow.Emit(list), func(app target.App_) (apps.App, bool) {
		return apps.App{App_: app}, true
	}), nil
}

func (a Apps) changes(ctx context.Context) (out <-chan apps.App, err error) {
	watch, err := fs2.NotifyWatch(ctx, a.Dir.Abs(), fsnotify.Create|fsnotify.Remove)
	if err != nil {
		return
	}
	return flow.Map(watch, func(event fsnotify.Event) (app apps.App, ok bool) {
		name := filepath.Base(event.Name)
		if _, err := object.ParseID(name); err != nil {
			return
		}
		switch event.Op {
		case fsnotify.Create:
			ta, err := a.Get(ctx, name)
			if err != nil {
				return
			}
			app.App_ = ta
		case fsnotify.Remove:
			open, err := a.Dir.Files().Open(filepath.Join("removed", name))
			if err != nil {
				return
			}
			defer open.Close()
			if err := json.NewDecoder(open).Decode(app.Manifest()); err != nil {
				return
			}
		}
		app.Event = &event
		ok = true
		return
	}), nil
}

func (a Apps) Uninstall(ctx context.Context, pkg string) error {
	app, err := a.Get(ctx, pkg)
	if err != nil {
		panic(err)
		return err
	}
	removed := filepath.Join(a.Dir.Abs(), "removed")
	if err = os.MkdirAll(removed, 0755); err != nil {
		panic(err)
		return err
	}
	removed = filepath.Join(removed, app.Path())
	removed = filepath.Join(removed, filepath.Base(app.Abs()))
	file, err := os.Create(removed)
	if err != nil {
		panic(err)
		return err
	}
	defer file.Close()
	if err = json.NewEncoder(file).Encode(app.Manifest()); err != nil {
		panic(err)
		return err
	}
	return os.Remove(app.Abs())
}

func (a Apps) Install(_ context.Context, reader io.ReadCloser) (err error) {
	defer reader.Close()
	tmpName := filepath.Join(a.Dir.Abs(), uuid.New().String())
	file, err := os.Create(tmpName)
	if err != nil {
		panic(err)
		return err
	}
	defer func() {
		_ = file.Close()
		if err != nil {
			_ = os.Remove(tmpName)
		}
	}()
	if _, err = io.Copy(file, reader); err != nil {
		panic(err)
		return
	}
	id, err := object.ResolveFile(tmpName)
	if err != nil {
		panic(err)
		return
	}
	name := filepath.Join(a.Dir.Abs(), id.String())
	if err = os.Rename(tmpName, name); err != nil {
		panic(err)
		return
	}
	return
}

func (a Apps) InstallSources(ctx context.Context, sources ...target.Source) (err error) {
	for _, app := range a.Resolve.List(sources...) {
		if err = a.Install(ctx, app.File()); err != nil {
			return
		}
	}
	return
}

func (a Apps) InstallFromPath(ctx context.Context, path string) error {
	found, err := a.Find(ctx, path)
	if err != nil {
		return err
	}
	for _, bundle := range found {
		if err = a.Install(ctx, bundle.Package().File()); err != nil {
			return err
		}
	}
	return nil

}
