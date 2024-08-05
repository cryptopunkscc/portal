package test

import (
	"context"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/find"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test__Target_Find(t *testing.T) {
	ctx := context.TODO()

	resolver := sources.Resolver[target.Portal_]()
	f := find.Combine(
		find.ById(resolver),
		find.ByPath(source.File, resolver)).
		Cached(&target.Cache[target.Portal_]{}).
		Reduced(
			target.Match[target.Project_],
			target.Match[target.Dist_],
			target.Match[target.Bundle_])

	found, err := f(ctx, "test.project.go")
	assert.ErrorIs(t, err, find.ErrNothing)

	if found, err = f(ctx, ""); err != nil {
		t.Fatal(err)
	}

	for i, base := range found {
		log.Println(i, base.Manifest().Name, base.Manifest().Package, reflect.TypeOf(base), base.Abs())
	}

	found, err = f(ctx, "test.project.go")
	assert.ErrorIs(t, err, nil)
}
