package test

import (
	"context"
	"github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/find"
	"github.com/cryptopunkscc/portal/target/portals"
	"testing"
)

func Test__Builder_find(t *testing.T) {
	f := find.New[target.Portal](deps{})

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

type deps struct{}

func (d deps) Path() target.Path                          { return apps.Path }
func (d deps) TargetFinder() target.Finder[target.Portal] { return portals.NewFind[target.Portal] }
func (d deps) TargetCache() *target.Cache[target.Portal]  { return target.NewCache[target.Portal]() }
