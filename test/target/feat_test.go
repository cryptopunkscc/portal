package test

import (
	"context"
	"github.com/cryptopunkscc/portal/feat"
	"github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/portals"
	"testing"
)

func Test__Builder_find(t *testing.T) {
	scope := feat.Scope[target.Portal]{
		GetPath:      apps.Path,
		TargetFinder: portals.NewFind[target.Portal],
		TargetCache:  target.NewCache[target.Portal](),
	}
	find := scope.GetTargetFind()

	for _, test := range portalTestCases {
		test := test
		t.Run(test.Src, func(t *testing.T) {
			apps_, err := find(context.TODO(), test.Src)
			if err != nil {
				t.Fatal(err)
			}
			for _, app := range apps_ {
				test.Assert(t, app)
			}
		})
	}
}
