package apps

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/source"
	"log"
	"reflect"
	"testing"
)

func TestResolveAll_List(t *testing.T) {
	file, err := source.File("../")
	if err != nil {
		t.Error(err)
	}

	for i, base := range target.List(ResolveAll, file) {
		log.Println(i, reflect.TypeOf(base), base.Abs())
	}

}

func TestResolveAll_Set(t *testing.T) {
	file, err := source.File("../")
	if err != nil {
		t.Error(err)
	}

	for i, base := range target.Set(ResolveAll, file) {
		log.Println(i, reflect.TypeOf(base), base.Abs())
	}

}
