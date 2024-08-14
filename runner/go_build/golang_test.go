package go_build

import (
	"context"
	golang "github.com/cryptopunkscc/portal/resolve/go"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoRunner_Run(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedGo)
	ctx := context.Background()
	project, err := golang.ResolveProject(src)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedGoBuild, project.Build())
	run := Runner()
	if err = run(ctx, project); err != nil {
		t.Fatal(err)
	}
}
