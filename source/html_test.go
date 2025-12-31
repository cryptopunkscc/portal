package source

import (
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var testHtmlMetadata = Manifest{
	Name:        "Name",
	Title:       "Title",
	Description: "Description",
	Package:     "Package",
	Version:     1,
	Icon:        "Icon",
	Runtime:     "Runtime",
	Type:        "Type",
}

var testIndexHtml = []byte(`portal.log("Hello Astral!!!")`)

func testHtmlRef(t *testing.T, name string) (ref Ref) {
	ref.Fs = afero.NewMemMapFs()
	//ref = Ref{afero.NewOsFs(), test.CleanMkdir(t, name)}
	return
}

func TestHtmlApp_WriteFs_ReadFS(t *testing.T) {
	ref := testHtmlRef(t, ".html_test_app")
	actual := HtmlApp{}
	expected := HtmlApp{}
	expected.Metadata = Metadata{Manifest: testHtmlMetadata}
	expected.IndexHtml = testIndexHtml
	test.NoError(t, expected.WriteRef(ref))
	test.NoError(t, actual.ReadSrc(&ref))

	expected.Fs = nil
	actual.Fs = nil
	actual.Path = ""
	require.Equal(t, expected, actual)
}

func TestHtmlProject_WriteFs_ReadFS(t *testing.T) {
	ref := testHtmlRef(t, ".html_test_project")
	actual := HtmlProject{}
	expected := HtmlProject{}
	expected.Manifest = ProjectMetadata{Manifest: testHtmlMetadata}
	expected.Html.IndexHtml = testIndexHtml
	test.NoError(t, expected.WriteRef(ref))
	test.NoError(t, actual.ReadSrc(&ref))

	expected.Fs = nil
	actual.Fs = nil
	actual.Path = ""
	require.Equal(t, expected, actual)
}

func TestHtmlBundle_WriteFs_ReadFS(t *testing.T) {
	ref := testHtmlRef(t, ".html_test_bundle")
	actual := HtmlBundle{}
	expected := HtmlBundle{}
	expected.App.Metadata = Metadata{Manifest: testHtmlMetadata}
	expected.Html.IndexHtml = testIndexHtml
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
