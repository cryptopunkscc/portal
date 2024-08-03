package sources

import (
	"github.com/cryptopunkscc/portal/apps"
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
	embed := source.Embed(apps.LauncherSvelteFS)

	for i, base := range target.List(Resolver[target.Base](), embed, file) {
		log.Println(i, reflect.TypeOf(base), base.Manifest().Package, base.Abs())
	}
}

func TestResolveAll_Set(t *testing.T) {
	file, err := source.File("../")
	if err != nil {
		t.Error(err)
	}

	for i, base := range target.Set(Resolver[target.Base](), file) {
		log.Println(i, reflect.TypeOf(base), base.Abs())
	}
}
