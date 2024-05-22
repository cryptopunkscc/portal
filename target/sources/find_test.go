package sources

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/bundle"
	"github.com/cryptopunkscc/go-astral-js/target/dist"
	js "github.com/cryptopunkscc/go-astral-js/target/js/embed"
	"github.com/cryptopunkscc/go-astral-js/target/npm"
	"github.com/cryptopunkscc/go-astral-js/target/project"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test_FromPath(t *testing.T) {
	assets := target.Abs("test_assets")
	targets := array.FromChan(FromPath[target.Portal](assets))

	for _, s := range targets {
		PrintTarget(s)
	}

	assert.Equal(t, 13, len(targets))
}

func Test_FindLibsInFs(t *testing.T) {
	targets := array.FromChan(FromFS[target.Source](js.PortalLibFS))

	for _, s := range targets {
		PrintTarget(s)
	}
}

func PrintTarget(t target.Source) {
	log.Println(reflect.TypeOf(t), t.Path(), t.Abs())
}

func Test_CustomFind(t *testing.T) {
	src := source.FromPath("test_assets")

	var find = target.Any[target.Source](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Lift(target.Try(npm.ResolveNodeModule))(
			target.Try(project.Resolve)),
		target.Try(dist.Resolve),
	)

	for s := range source.Stream[target.Source](find, src) {
		PrintTarget(s)
	}
}
