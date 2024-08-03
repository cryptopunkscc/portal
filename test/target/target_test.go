package test

import (
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/target"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test__target_Any__test_assets(t *testing.T) {
	var resolve = sources.Resolver[target.Portal_]()
	src, err := source.File("test_data")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range target.List(resolve, src) {
		PrintTarget(s)
	}
}

func PrintTarget(t target.Source) {
	log.Println(reflect.TypeOf(t), t.Path(), t.Abs())
}

func Test__target_Type(t *testing.T) {
	assert.True(t, target.TypeFrontend.Is(target.TypeAny))
}
