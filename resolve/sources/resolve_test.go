package sources

import (
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"reflect"
	"testing"
)

func TestResolveAll_List(t *testing.T) {
	file, err := source.File("../")
	if err != nil {
		t.Error(err)
	}
	embed := source.Embed(apps.LauncherSvelteFS)

	for i, base := range target.List(Resolver[target.Portal_](), embed, file) {
		log.Println(i, reflect.TypeOf(base), base.Manifest().Package, base.Abs())
	}
}

func TestResolveAll_Set(t *testing.T) {
	file, err := source.File("../")
	if err != nil {
		t.Error(err)
	}

	for i, base := range target.Set(Resolver[target.Portal_](), file) {
		log.Println(i, reflect.TypeOf(base), base.Abs())
	}
}
