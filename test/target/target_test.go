package test

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/project"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test__target_Any__test_assets(t *testing.T) {
	var find = target.Any[target.Project](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Try(npm.ResolveNodeModule).Lift(
			target.Try(project.ResolveNpm)),
		target.Try(project.ResolveGo),
		target.Try(dist.Resolve),
	)
	src := source.FromPath("test_data")
	for _, s := range source.List[target.Project](find, src) {
		PrintTarget(s)
	}
}

func PrintTarget(t target.Source) {
	log.Println(reflect.TypeOf(t), t.Path(), t.Abs())
}

func Test__target_Type(t *testing.T) {
	assert.True(t, target.TypeFrontend.Is(target.TypeAny))
}
