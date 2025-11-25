package test

import (
	"os"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/assets"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/project"
)

func TestCreateProject(t *testing.T) {
	name := "created_project"
	_ = os.RemoveAll(name)

	opts := project.CreateOpts{}
	opts.Path = name
	opts.TemplatesFS = assets.ArrayFs{}
	err := project.Create(opts)
	test.AssertErr(t, err)
}
