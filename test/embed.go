package test

import (
	"embed"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/zip"
	"log"
)

//go:embed data.zip
var EmbedFS embed.FS
var EmbedRoot target.Source = source.Embed(EmbedFS)
var (
	EmbedSh             target.Source
	EmbedJs             target.Source
	EmbedHtml           target.Source
	EmbedJsRollup       target.Source
	EmbedJsRollupDist   target.Source
	EmbedJsRollupBundle target.Source
	EmbedSvelte         target.Source
	EmbedSvelteDist     target.Source
	EmbedSvelteBundle   target.Source
	EmbedBundle         target.Source
	EmbedGo             target.Source
	EmbedGoDist         target.Source
	EmbedGoBundle       target.Source
)

func init() {
	var err error
	if EmbedRoot, err = EmbedRoot.Sub("data.zip"); err != nil {
		panic(err)
	}
	if EmbedRoot, err = zip.Resolve(EmbedRoot); err != nil {
		panic(err)
	}
	for path, src := range map[string]*target.Source{
		"sh":                   &EmbedSh,
		"js":                   &EmbedJs,
		"html":                 &EmbedHtml,
		"bundle/bundle.portal": &EmbedBundle,

		"js-rollup":      &EmbedJsRollup,
		"js-rollup/dist": &EmbedJsRollupDist,
		"js-rollup/build/new.portal.js-rollup_.portal": &EmbedJsRollupBundle,

		"svelte":                                 &EmbedSvelte,
		"svelte/dist":                            &EmbedSvelteDist,
		"svelte/build/new.portal.svelte_.portal": &EmbedSvelteBundle,

		"go":                       &EmbedGo,
		"go/dist":                  &EmbedGoDist,
		"go/build/test.go_.portal": &EmbedGoBundle,
	} {
		if *src, err = EmbedRoot.Sub(path); err != nil {
			log.Printf("%s: %v", path, err)
		}
	}
}

var (
	EmbedShManifest        = &target.Manifest{Name: "sh", Package: "new.portal.sh", Exec: "main"}
	EmbedJsManifest        = &target.Manifest{Name: "js", Package: "new.portal.js"}
	EmbedBundleManifest    = &target.Manifest{Name: "bundle", Package: "new.portal.bundle"}
	EmbedHtmlManifest      = &target.Manifest{Name: "html", Package: "new.portal.html"}
	EmbedSvelteManifest    = &target.Manifest{Name: "svelte", Package: "new.portal.svelte"}
	EmbedJsRollupManifest  = &target.Manifest{Name: "js-rollup", Package: "new.portal.js-rollup"}
	EmbedGoManifest        = &target.Manifest{Name: "go", Package: "test.go", Title: "test go", Exec: "main"}
	EmbedGoProjectManifest = &target.Manifest{Name: "go", Package: "test.go", Title: "test go"}
)

var (
	EmbedSvelteBuild   = target.Builds(nil)
	EmbedJsRollupBuild = target.Builds{
		"default": target.Build{Cmd: "cmd2", Deps: []string{"dep2"}},
		"linux":   target.Build{Cmd: "cmd3", Deps: []string{"dep2", "dep3"}},
		"windows": target.Build{Cmd: "cmd4", Deps: []string{"dep2", "dep4"}},
	}
	EmbedGoBuild = target.Builds{
		"default": target.Build{Out: "main", Cmd: "go build -o dist/main"},
		"linux":   target.Build{Out: "main", Cmd: "go build -o dist/main", Deps: []string{"gcc", "libgtk-3-dev", "libayatana-appindicator3-dev"}},
		"windows": target.Build{Out: "main.exe", Cmd: "go build -ldflags -H=windowsgui -o dist/main.exe", Env: []string{"CGO_ENABLED=1"}},
	}
)
