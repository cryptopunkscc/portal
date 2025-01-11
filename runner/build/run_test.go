package build

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"testing"
)

func TestFeat_Run(t *testing.T) {
	feat := NewRunner(
		func(s string) error {
			t.Logf("clean %s", s)
			return nil
		},
		func(ctx context.Context, src target.Project_, _ ...string) (err error) {
			t.Logf("run dist %T %s %s", src, src.Manifest().Package, src.Abs())
			return nil
		},
		func(ctx context.Context, src target.Dist_, _ ...string) (err error) {
			t.Logf("run pack %T %s %s", src, src.Manifest().Package, src.Abs())
			return nil
		},
	)
	err := feat.Run(context.TODO(), "../../example")
	if err != nil {
		t.Fatal(err)
	}
}
