package test

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/portals"
	"testing"
)

func Test__Builder_find(t *testing.T) {
	scope := feat.Scope[target.Portal]{
		GetPath:      apps.Path,
		TargetFinder: portals.NewFind,
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
