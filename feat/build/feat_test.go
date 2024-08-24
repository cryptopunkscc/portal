package build

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"testing"
)

func TestFeat_Run(t *testing.T) {
	feat := NewFeat(
		func(s string) error {
			t.Logf("clean %s", s)
			return nil
		},
		func(ctx context.Context, src target.Project_) (err error) {
			t.Logf("run dist %T %s %s", src, src.Manifest().Package, src.Abs())
			return nil
		},
		func(ctx context.Context, src target.Dist_) (err error) {
			t.Logf("run pack %T %s %s", src, src.Manifest().Package, src.Abs())
			return nil
		},
	)
	err := feat.Run(context.TODO(), "../../test/data")
	if err != nil {
		t.Fatal(err)
	}
}
