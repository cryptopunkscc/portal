package test

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/test/docker"
	"testing"
	"time"
)

func TestE2E_2(t *testing.T) {
	dc := docker.Container{
		Bin:           "podman",
		Image:         "e2e-test",
		Network:       "e2e-test-net",
		Logfile:       "portald.log",
		InstallerPath: "/portal/bin/install-portal-to-astral",
	}
	var c = []*Cases{
		{id: 0, Device: dc.New(0)},
		{id: 1, Device: dc.New(1)},
	}

	runner := test.Runner{}
	tests := []test.Task{
		// ====== base ======
		{
			Name: "build",
			Test: dc.BuildImage(),
		},
		{
			Name: "print install help",
			Test: c[0].PrintInstallHelp(),
		},
		{
			Name: "start portald via portal",
			Test: c[0].PortalStart(),
		},
		{
			Name: "start portald via portal",
			Test: c[1].PortalStartAwait(),
		},
		{
			Name: "portal help",
			Test: c[0].PortalHelp(),
		},
		{
			Name: "user claim",
			Test: c[0].UserClaim(c[1]),
		},
		// ====== dev.js ======
		{
			Name: "list js templates",
			Test: c[0].ListTemplates("dev.js"),
		},
		{
			Name: "create js project",
			Test: c[0].NewProject(jsProject),
		},
		{
			Name: "pack js project",
			Test: c[0].BuildProject(jsProject),
		},
		{
			Name: "publish js project",
			Test: c[0].PublishProject(jsProject),
		},
		{
			Name: "list available apps",
			Test: c[1].ListAvailableApps(jsProject),
			Require: test.Tests{
				c[0].UserClaim(c[1]),
				c[0].PublishProject(jsProject),
			},
		},
		{
			Name: "install available app by name",
			Test: c[1].InstallAvailableApp(jsProject),
			Require: test.Tests{
				c[0].UserClaim(c[1]),
				c[0].PublishProject(jsProject),
			},
		},
		{
			Name:    "run js app",
			Test:    c[1].RunApp(jsProject),
			Require: test.Tests{c[1].InstallAvailableApp(jsProject)},
		},
		{
			Name: "create js-rollup project",
			Test: c[0].NewProject(jsRollupProject),
		},
		{
			Name: "build js-rollup project",
			Test: c[0].BuildProject(jsRollupProject),
		},
		{
			Name: "publish js-rollup app",
			Test: c[0].PublishProject(jsRollupProject),
		},
		{
			Name: "install js-rollup app by name",
			Test: c[1].InstallAvailableApp(jsRollupProject),
			Require: test.Tests{
				c[0].UserClaim(c[1]),
				c[0].PublishProject(jsRollupProject),
			},
		},
		{
			Name:    "run js-rollup app",
			Test:    c[1].RunApp(jsRollupProject),
			Require: test.Tests{c[1].InstallAvailableApp(jsRollupProject)},
		},
		// ====== dev.html ======
		{
			Name: "list html templates",
			Test: c[0].ListTemplates("dev.html"),
		},
		{
			Name: "create html project",
			Test: c[0].NewProject(htmlProject),
		},
		{
			Name: "create svelte project",
			Test: c[0].NewProject(svelteProject),
		},
		{
			Name: "build svelte project",
			Test: c[0].BuildProject(svelteProject),
		},
		{
			Name: "create react project",
			Test: c[0].NewProject(reactProject),
		},
		{
			Name: "build react project",
			Test: c[0].BuildProject(reactProject),
		},
		// ====== tear down ======
		{
			Name:  "portal close",
			Test:  c[0].PortalClose(),
			Group: 1,
		},
		{
			Name:  "print logs",
			Group: 2,
			Test:  c[0].PrintLogs(),
		},
		{
			Name:  "print logs",
			Group: 3,
			Test:  c[1].PrintLogs(),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d  %s", i, tt.Name), runner.Run(tests, tt))
	}

	t.Cleanup(func() {
		time.Sleep(1 * time.Second) // await logs
		for _, ccc := range c {
			ccc.Stop()
		}
	})
}

var jsProject = ProjectOpts{
	runner:   "dev.js",
	template: "js",
	name:     "js-raw",
}

var jsRollupProject = ProjectOpts{
	runner:   "dev.js",
	template: "js-rollup",
}

var htmlProject = ProjectOpts{
	runner:   "dev.html",
	template: "html",
}

var svelteProject = ProjectOpts{
	runner:   "dev.html",
	template: "svelte",
}

var reactProject = ProjectOpts{
	runner:   "dev.html",
	template: "react",
}
