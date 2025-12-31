package source

import (
	"testing"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var testMetadata = Manifest{
	Name:        "Name",
	Title:       "Title",
	Description: "Description",
	Package:     "Package",
	Version:     1,
	Icon:        "Icon",
	Runtime:     "Runtime",
	Type:        "Type",
}

var testMainJs = []byte(`portal.log("Hello Astral!!!")`)

func testDir(t *testing.T, name string) (dir afero.Fs) {
	dir = afero.NewMemMapFs()
	dir = afero.NewBasePathFs(
		afero.NewOsFs(),
		test.CleanMkdir(t, name),
	)
	return
}

func TestJsDist_WriteFs_ReadFS(t *testing.T) {
	dir := testDir(t, ".js_test_dist")
	actual := JsDist{}
	expected := JsDist{}
	expected.Metadata = Metadata{Manifest: testMetadata}
	expected.MainJs = testMainJs
	require.NoError(t, expected.WriteFs(dir))
	require.NoError(t, actual.ReadFs(dir))

	expected.Fs = nil
	actual.Fs = nil
	require.Equal(t, expected, actual)
}

func TestJsProject_WriteFs_ReadFS(t *testing.T) {
	dir := testDir(t, ".js_test_project")
	actual := JsProject{}
	expected := JsProject{}
	expected.Manifest = ProjectMetadata{Manifest: testMetadata}
	expected.Js.MainJs = testMainJs
	require.NoError(t, expected.WriteFs(dir))
	require.NoError(t, actual.ReadFs(dir))

	expected.Fs = nil
	actual.Fs = nil
	require.Equal(t, expected, actual)
}

func TestJsBundle_WriteFs_ReadFS(t *testing.T) {
	dir := testDir(t, ".js_test_bundle")
	actual := JsBundle{}
	expected := JsBundle{}
	expected.Dist.Metadata = Metadata{Manifest: testMetadata}
	expected.Js.MainJs = testMainJs
	err := expected.WriteFs(dir)
	plog.Println(err)
	require.NoError(t, err)

	actual.Name = expected.Name
	require.NoError(t, actual.ReadFs(dir))

	expected.Fs = nil
	expected.Dist.Fs = nil
	expected.ZipFs = nil
	actual.Fs = nil
	actual.Dist.Fs = nil
	actual.ZipFs = nil
	require.Equal(t, expected, actual)
}
