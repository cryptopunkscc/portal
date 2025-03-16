package build_all

import (
	"context"
	"testing"
)

func TestBuild(t *testing.T) {
	if err := Run(context.Background(), "../../apps"); err != nil {
		t.Fatal(err)
	}
}
