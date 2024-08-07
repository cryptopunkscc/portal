package sources

import (
	"context"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
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
	"/test_data/go":                       test.EmbedGoManifest,
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
	"/test_data/go":                   test.EmbedGoManifest,
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

func path(src string) (path string, err error) {
	path, ok := packagePathReduced[src]
	if !ok {
		err = target.ErrNotTarget
	}
	return

}

var findByPath = target.FindByPath(
	source.File,
	Resolver[target.Portal_](),
)

var findByPathReduced = findByPath.Reduced(
	target.Match[target.Project_],
	target.Match[target.Bundle_],
	target.Match[target.Dist_],
)

func TestByPath(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedRoot)
	expected := pathManifestAll
	find := target.FindByPath(
		source.File,
		Resolver[target.Portal_](),
	)

	portals := test.Assert(find(context.Background(), src.Abs()))

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
	find := target.FindByPath(
		source.File,
		Resolver[target.Portal_](),
	).Reduced(
		target.Match[target.Project_],
		target.Match[target.Bundle_],
		target.Match[target.Dist_],
	)

	portals := test.Assert(find(context.Background(), src.Abs()))

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
	find := findByPathReduced.Cached(&target.Cache[target.Portal_]{})

	test.Assert(find(context.Background(), src.Abs()))

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
	ctx := context.Background()
	expected := packageManifestReduced
	find := findByPathReduced.ById(path)
	for pkg, manifest := range expected {
		t.Log(manifest)
		r := test.Assert(find(ctx, pkg))[0]
		assert.Equal(t, manifest, expected[pkg], r.Manifest())
		delete(expected, pkg)
	}
	assert.Equal(t, make(map[string]*target.Manifest), expected)
}
