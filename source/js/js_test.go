package js

import (
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var testJsManifest = app.Manifest{
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

func testJsRef(t *testing.T, name string) (ref source.Ref) {
	ref.Fs = afero.NewMemMapFs()
	//ref = Ref{Fs: afero.NewOsFs(), Path: test.CleanMkdir(t, name)}
	return
}

func TestJsApp_WriteFs_ReadFS(t *testing.T) {
	ref := testJsRef(t, ".js_test_app")
	actual := App{}
	expected := App{}
	expected.Metadata = app.Metadata{Manifest: testJsManifest}
	expected.MainJs = testMainJs
	test.NoError(t, expected.WriteRef(ref))
	test.NoError(t, actual.ReadSrc(&ref))

	expected.Fs = nil
	actual.Fs = nil
	actual.Path = ""
	require.Equal(t, expected, actual)
}

func TestJsProject_WriteFs_ReadFS(t *testing.T) {
	ref := testJsRef(t, ".js_test_project")
	actual := Project{}
	expected := Project{}
	expected.Manifest = testJsManifest
	expected.Js.MainJs = testMainJs
	test.NoError(t, expected.WriteRef(ref))
	test.NoError(t, actual.ReadSrc(&ref))

	expected.Fs = nil
	actual.Fs = nil
	actual.Path = ""
	require.Equal(t, expected, actual)
}

func TestJsBundle_WriteFs_ReadFS(t *testing.T) {
	ref := testJsRef(t, ".js_test_bundle")
	actual := JsBundle{}
	expected := JsBundle{}
	expected.Dist.Metadata = app.Metadata{Manifest: testJsManifest}
	expected.Js.MainJs = testMainJs
	test.NoError(t, expected.WriteRef(ref))

	ref.Path = expected.Zip.File.Path
	test.NoError(t, actual.ReadSrc(&ref))

	expected.Fs = nil
	expected.Unpacked = nil
	actual.Fs = nil
	actual.Unpacked = nil
	actual.Blob = nil
	require.Equal(t, expected, actual)
}
