package build

import (
	"context"
	"testing"
)

func TestBuild(t *testing.T) {
	if err := Create().Run(context.Background(), "../../apps"); err != nil {
		t.Fatal(err)
	}
}
