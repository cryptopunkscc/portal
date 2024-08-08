package srv

import (
	"context"
	"github.com/cryptopunkscc/portal/resolve/sources"
	. "github.com/cryptopunkscc/portal/target"
	"testing"
)

func TestModule_TargetFind(t *testing.T) {
	m := &module[Portal_]{}
	m.Deps = m
	portals, err := m.TargetFind().Call(context.Background(), "../../apps")
	if err != nil {
		t.Fatal(err)
	}
	for i, portal := range portals {
		t.Logf("%d %T %s %s", i, portal, portal.Manifest().Package, portal.Abs())
	}
}

type module[T Portal_] struct{ Module[T] }

func (m *module[T]) TargetResolve() Resolve[T] { return sources.Resolver[T]() }
func (m *module[T]) Priority() Priority {
	return Priority{
		Match[Project_],
		Match[Dist_],
		Match[Bundle_],
	}
}
