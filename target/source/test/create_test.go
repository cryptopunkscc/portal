package test

import (
	"embed"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

//go:embed template
var templateFS embed.FS

func TestCreate(t *testing.T) {
	opts := source.CreateOpts{}
	opts.Path = "created_source"
	opts.Template = "template"
	opts.TemplateArgs = source.TemplateArgs{"Bar": "baz"}
	opts.TemplatesFS = templateFS

	_ = os.RemoveAll(opts.Path)
	err := source.Create(opts)
	test.AssertErr(t, err)
	assert.FileExists(t, filepath.Join(opts.Path, "foo.yml"))
}
