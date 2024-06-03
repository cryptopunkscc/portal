package test

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/bundle"
	"github.com/cryptopunkscc/go-astral-js/target/dist"
	"github.com/cryptopunkscc/go-astral-js/target/npm"
	"github.com/cryptopunkscc/go-astral-js/target/project"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test__target_Any__test_assets(t *testing.T) {
	var find = target.Any[target.Source](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Try(npm.ResolveNodeModule).Lift(
			target.Try(project.Resolve)),
		target.Try(dist.Resolve),
	)
	src := source.FromPath("test_assets")
	for _, s := range source.List[target.Source](find, src) {
		PrintTarget(s)
	}
}

func PrintTarget(t target.Source) {
	log.Println(reflect.TypeOf(t), t.Path(), t.Abs())
}

func Test__target_Type(t *testing.T) {
	assert.True(t, target.TypeFrontend.Is(target.TypeAny))
}
