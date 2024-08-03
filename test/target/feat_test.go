package test

import (
	"context"
	apps2 "github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/find"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/source"
	"github.com/cryptopunkscc/portal/target2/sources"
	"testing"
)

func Test__Builder_find(t *testing.T) {
	f := find.Inject[target.Base](&deps{})

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

type deps struct {
	cache target.Cache[target.Base]
}

func (d *deps) Priority() target.Priority {
	return target.Priority{
		target.Match[target.Project_],
		target.Match[target.Bundle_],
		target.Match[target.Dist_],
	}
}
func (d *deps) Path() target.Path                          { return apps.Path }
func (d *deps) Embed() []target.Source                     { return []target.Source{source.Embed(apps2.LauncherSvelteFS)} }
func (d *deps) TargetFile() target.File                    { return source.File }
func (d *deps) TargetCache() *target.Cache[target.Base]    { return &d.cache }
func (d *deps) TargetResolve() target.Resolve[target.Base] { return sources.Resolver[target.Base]() }
