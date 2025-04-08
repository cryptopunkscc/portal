package sources

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

var wd string

func init() {
	var err error
	if wd, err = os.Getwd(); err != nil {
		panic(err)
	}
}

var pathManifestAll = map[string]*target.Manifest{
	"/test_data/bundle/bundle.portal":     test.EmbedBundleManifest,
	"/test_data/go":                       test.EmbedGoProjectManifest,
	"/test_data/go/build/test.go_.portal": test.EmbedGoManifest,
	"/test_data/go/dist":                  test.EmbedGoManifest,
	"/test_data/html":                     test.EmbedHtmlManifest,
	"/test_data/js":                       test.EmbedJsManifest,
	"/test_data/js-rollup":                test.EmbedJsRollupManifest,
	"/test_data/js-rollup/build/new.portal.js-rollup_.portal": test.EmbedJsRollupManifest,
	"/test_data/js-rollup/dist":                               test.EmbedJsRollupManifest,
	"/test_data/sh":                                           test.EmbedShManifest,
	"/test_data/svelte":                                       test.EmbedSvelteManifest,
	"/test_data/svelte/build/new.portal.svelte_.portal":       test.EmbedSvelteManifest,
	"/test_data/svelte/dist":                                  test.EmbedSvelteManifest,
}

var pathManifestReduced = map[string]*target.Manifest{
	"/test_data/go":                   test.EmbedGoProjectManifest,
	"/test_data/svelte":               test.EmbedSvelteManifest,
	"/test_data/js-rollup":            test.EmbedJsRollupManifest,
	"/test_data/bundle/bundle.portal": test.EmbedBundleManifest,
	"/test_data/sh":                   test.EmbedShManifest,
	"/test_data/js":                   test.EmbedJsManifest,
	"/test_data/html":                 test.EmbedHtmlManifest,
}

var packageManifestReduced = map[string]*target.Manifest{
	test.EmbedGoManifest.Package:       test.EmbedGoManifest,
	test.EmbedSvelteManifest.Package:   test.EmbedSvelteManifest,
	test.EmbedJsRollupManifest.Package: test.EmbedJsRollupManifest,
	test.EmbedBundleManifest.Package:   test.EmbedBundleManifest,
	test.EmbedShManifest.Package:       test.EmbedShManifest,
	test.EmbedJsManifest.Package:       test.EmbedJsManifest,
	test.EmbedHtmlManifest.Package:     test.EmbedHtmlManifest,
}

var packagePathReduced = map[string]string{
	test.EmbedGoManifest.Package:       "test_data/go",
	test.EmbedSvelteManifest.Package:   "test_data/svelte",
	test.EmbedJsRollupManifest.Package: "test_data/js-rollup",
	test.EmbedBundleManifest.Package:   "test_data/bundle/bundle.portal",
	test.EmbedShManifest.Package:       "test_data/sh",
	test.EmbedJsManifest.Package:       "test_data/js",
	test.EmbedHtmlManifest.Package:     "test_data/html",
}

var provider = target.Provider[target.Portal_]{
	Priority: target.Priority{
		target.Match[target.Project_],
		target.Match[target.Bundle_],
		target.Match[target.Dist_],
	},
	Repository: source.Repository,
	Resolve:    Resolver[target.Portal_](),
}

func TestByPath(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedRoot)
	expected := pathManifestAll

	portals := provider.All(src.Abs())

	for _, portal := range portals {
		key := strings.TrimPrefix(portal.Abs(), wd)
		assert.Equal(t, expected[key], portal.Manifest(), target.Sprint(portal))
		delete(expected, key)
	}
	assert.Equal(t, make(map[string]*target.Manifest), expected)
}

func TestByPathReduced(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedRoot)
	expected := pathManifestReduced

	portals := provider.Provide(src.Abs())

	for _, portal := range portals {
		key := strings.TrimPrefix(portal.Abs(), wd)
		assert.Equal(t, expected[key], portal.Manifest(), target.Sprint(portal))
		delete(expected, key)
	}
	assert.Equal(t, make(map[string]*target.Manifest), expected)
}

func Test_ByPath_Reduced_Cached(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedRoot)
	expected := packageManifestReduced

	assert.NotNil(t, provider.Provide(src.Abs()))

	for pkg, manifest := range expected {
		t.Log(manifest)
		assert.Equal(t, manifest, expected[pkg], pkg)
		delete(expected, pkg)
	}
	assert.Equal(t, make(map[string]*target.Manifest), expected)
}

func Test_ByPath_Reduced_ById(t *testing.T) {
	defer test.Clean()
	test.Copy(test.EmbedRoot)
	expected := packageManifestReduced
	for pkg, manifest := range expected {
		t.Log(manifest)
		p := provider.Provide(pkg)
		assert.NotNil(t, p)
		r := p[0]
		assert.Equal(t, manifest, expected[pkg], r.Manifest())
		delete(expected, pkg)
	}
	assert.Equal(t, make(map[string]*target.Manifest), expected)
}
