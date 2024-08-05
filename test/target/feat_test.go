package test

import (
	"context"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/find"
	"testing"
)

func Test__Builder_find(t *testing.T) {
	f := find.
		ByPath(source.File, sources.Resolver[target.Portal_]()).
		Cached(&target.Cache[target.Portal_]{}).
		Reduced(
			target.Match[target.Project_],
			target.Match[target.Bundle_],
			target.Match[target.Dist_])

	for _, test := range portalTestCases {
		test := test
		t.Run(test.Src, func(t *testing.T) {
			apps_, err := f(context.TODO(), test.Src)
			if err != nil {
				t.Fatal(err)
			}
			for _, app := range apps_ {
				test.Assert(t, app)
			}
		})
	}
}
