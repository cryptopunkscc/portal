package test

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var portalTestCases = []Case[string]{
	{Src: ".", Matchers: []*Target{
		RpcBackend,
		RpcFrontend,
		BasicBackend,
		BasicFrontend,
		ProjectFrontend,
		ProjectBackend,
		ProjectGo,
		Launcher,
		DistExecutable,
	}},
	{Src: "test_data/rpc", Matchers: []*Target{
		RpcBackend,
		RpcFrontend,
	}},
	{Matcher: Launcher, Src: Launcher.Manifest.Name},
	{Matcher: Launcher, Src: Launcher.Manifest.Package},
	{Matcher: BasicBackend, Src: BasicBackend.Abs},
	{Matcher: BasicFrontend, Src: BasicFrontend.Abs},
	{Matcher: RpcFrontend, Src: RpcFrontend.Abs},
	{Matcher: RpcBackend, Src: RpcBackend.Abs},
	{Matcher: ProjectBackend, Src: ProjectBackend.Abs},
	{Matcher: ProjectFrontend, Src: ProjectFrontend.Abs},
}

var Launcher = &Target{
	Path:     ".",
	Abs:      "launcher/svelte/dist",
	Manifest: &target.Manifest{Name: "launcher", Title: "Portal Launcher", Description: "Portal applications launcher.", Package: "cc.cryptopunks.portal.launcher", Version: "0.0.0", Icon: "icon.svg"},
}

var BasicBackend = &Target{
	Path:     ".",
	Abs:      "test_data/basic/back",
	Manifest: &target.Manifest{Name: "test-basic-back", Title: "Example basic backend", Description: "", Package: "test.basic.back", Version: "0.0.0", Icon: ""},
}

var BasicFrontend = &Target{
	Path:     ".",
	Abs:      "test_data/basic/front",
	Manifest: &target.Manifest{Name: "test-basic-front", Title: "test basic frontend", Description: "", Package: "test.basic.ui", Version: "0.0.0", Icon: ""},
}

var RpcFrontend = &Target{
	Path:     ".",
	Abs:      "test_data/rpc/front",
	Manifest: &target.Manifest{Name: "test-rpc-front", Title: "test rpc frontend", Description: "", Package: "test.rpc.front", Version: "0.0.0", Icon: ""},
}

var RpcFrontendBundle = &Target{
	Path:     "test-rpc-front_0.0.0.portal",
	Abs:      "test_data/rpc/build",
	Manifest: &target.Manifest{Name: "test-rpc-front", Title: "test rpc frontend", Description: "", Package: "test.rpc.front", Version: "0.0.0", Icon: ""},
}

var RpcBackend = &Target{
	Path:     ".",
	Abs:      "test_data/rpc/back",
	Manifest: &target.Manifest{Name: "test-rpc-back", Title: "test rpc backend", Description: "", Package: "test.rpc.back", Version: "0.0.0", Icon: ""},
}

var RpcBackendBundle = &Target{
	Path:     "test-rpc-back_0.0.0.portal",
	Abs:      "test_data/rpc/build",
	Manifest: &target.Manifest{Name: "test-rpc-back", Title: "test rpc backend", Description: "", Package: "test.rpc.back", Version: "0.0.0", Icon: "", Exec: ""},
}

var ProjectBackend = &Target{
	Path:     ".",
	Abs:      "test_data/project/backend",
	Manifest: &target.Manifest{Name: "test-project-backend", Title: "test project backend", Description: "", Package: "test.project.backend", Version: "0.0.0", Icon: ""},
}

var ProjectFrontend = &Target{
	Path:     ".",
	Abs:      "test_data/project/svelte",
	Manifest: &target.Manifest{Name: "test-project-svelte", Title: "test project svelte", Description: "", Package: "test.project.svelte", Version: "0.0.0", Icon: ""},
}

var ProjectGo = &Target{
	Path:     ".",
	Abs:      "test_data/project/go",
	Manifest: &target.Manifest{Name: "test-project-go", Title: "test project go", Description: "", Package: "test.project.go", Version: "0.0.0", Icon: ""},
}

var DistExecutable = &Target{
	Path:     ".",
	Abs:      "test_data/exec/sh",
	Manifest: &target.Manifest{Name: "test-dist-sh", Title: "test dist sh", Description: "", Package: "test.dist.sh", Version: "0.0.0", Icon: "", Exec: "exec.sh", Build: "", Env: target.Env{Timeout: 0}},
}

var BundleExecutable = &Target{
	Path:     "test-exec_0.0.0.portal",
	Abs:      "test_data/build",
	Manifest: &target.Manifest{Name: "test-exec", Title: "test exec", Description: "", Package: "test.exec", Version: "0.0.0", Icon: "", Exec: "exec.sh"},
}

type Case[T any] struct {
	Src      T
	Matcher  *Target
	Matchers []*Target
}

func (c Case[T]) Assert(t *testing.T, portal target.Portal_) {
	if c.Matcher != nil {
		c.Matcher.Assert(t, portal)
		return
	}
	for _, matcher := range c.Matchers {
		if matcher.Manifest.Package != portal.Manifest().Package {
			continue
		}
		matcher.Assert(t, portal)
		return
	}
	t.Error("no target matcher for:", portal.Abs(), portal.Manifest().Package)
}

type Target struct {
	Path     string
	Abs      string
	Manifest *target.Manifest
}

func (p Target) Assert(t *testing.T, portal target.Portal_) {
	assert.NotNil(t, portal)
	assert.Contains(t, portal.Abs(), p.Abs)
	assert.True(t, strings.HasSuffix(portal.Abs(), p.Abs), "\n%s\n%s", portal.Abs(), p.Abs)
	assert.Equal(t, p.Path, portal.Path())
	assert.Equal(t, p.Manifest, portal.Manifest())
}
