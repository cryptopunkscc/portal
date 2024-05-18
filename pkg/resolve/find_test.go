package resolve

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	js "github.com/cryptopunkscc/go-astral-js/pkg/binding/out"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test_FromPath(t *testing.T) {
	assets := target.Abs("test_assets")
	targets := array.FromChan(FromPath[target.Source](assets))

	for _, source := range targets {
		PrintTarget(source)
	}

	assert.Equal(t, 8, len(targets))
}

func Test_FindLibsInFs(t *testing.T) {
	targets := array.FromChan(FromFS[target.Source](js.PortalLibFS))

	for _, source := range targets {
		PrintTarget(source)
	}
}

func PrintTarget(t target.Source) {
	log.Println(reflect.TypeOf(t), t.Path(), t.Abs())
}
