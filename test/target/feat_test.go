package test

import (
	"context"
	apps2 "github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/find"
	npm2 "github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"testing"
)

func Test__Builder_find(t *testing.T) {
	f := find.Inject[target.Portal_](&deps{})

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
	cache target.Cache[target.Portal_]
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
func (d *deps) TargetCache() *target.Cache[target.Portal_] { return &d.cache }
func (d *deps) TargetResolve() target.Resolve[target.Portal_] {
	return sources.Resolver[target.Portal_]()
}

func TestCase_Assert(t *testing.T) {
	//_ = fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, err error) error {
	//	if filepath.Base(path) == "node_modules" {
	//		return fs.SkipDir
	//	}
	//	log.Println(path)
	//	return err
	//})

	src, _ := source.File(".")
	for _, p := range target.List(
		target.Any[target.NodeModule](
			target.Skip("node_modules"),
			target.Try(npm2.Resolve),
		),
		src,
	) {
		log.Println(p.Abs())
	}
}
