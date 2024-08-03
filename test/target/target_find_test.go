package test

import (
	"context"
	apps2 "github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/find"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/target"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test__Target_Find(t *testing.T) {
	ctx := context.TODO()
	f := find.New[target.Portal_](
		&target.Cache[target.Portal_]{},
		apps.Path,
		source.File,
		sources.Resolver[target.Portal_](),
		target.Priority{
			target.Match[target.Project_],
			target.Match[target.Dist_],
			target.Match[target.Bundle_],
		},
		source.Embed(apps2.LauncherSvelteFS),
	)

	found, err := f(ctx, "test.project.go")
	assert.ErrorIs(t, err, find.ErrNoPortals)

	if found, err = f(ctx, ""); err != nil {
		t.Fatal(err)
	}

	for i, base := range found {
		log.Println(i, base.Manifest().Name, base.Manifest().Package, reflect.TypeOf(base), base.Abs())
	}

	found, err = f(ctx, "test.project.go")
	assert.ErrorIs(t, err, nil)

	found, err = f(ctx, "cc.cryptopunks.portal.launcher")
	assert.ErrorIs(t, err, nil)
}
