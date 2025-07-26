package test

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/test"
	"testing"
	"time"
)

func TestE2E_2(t *testing.T) {
	c := create(2, container{
		image:   "e2e-test",
		network: "e2e-test-net",
		logfile: "portald.log",
	})

	runner := test.Runner{}
	tests := []test.Task{
		// ====== base ======
		{
			Name: "build",
			Test: c[0].buildImage(),
		},
		{
			Name: "print install help",
			Test: c[0].printInstallHelp(),
		},
		{
			Name: "start portald via portal",
			Test: c[0].portalStart(),
		},
		{
			Name: "start portald via portal",
			Test: c[1].portalStartAwait(),
		},
		{
			Name: "portal help",
			Test: c[0].portalHelp(),
		},
		{
			Name: "user claim",
			Test: c[0].userClaim(c[1]),
		},
		// ====== dev.js ======
		{
			Name: "list js templates",
			Test: c[0].listTemplates("dev.js"),
		},
		{
			Name: "create js project",
			Test: c[0].newProject(jsProject),
		},
		{
			Name: "pack js project",
			Test: c[0].buildProject(jsProject),
		},
		{
			Name: "publish js project",
			Test: c[0].publishProject(jsProject),
		},
		{
			Name: "list available apps",
			Test: c[1].listAvailableApps(jsProject),
			Require: test.Tests{
				c[0].userClaim(c[1]),
				c[0].publishProject(jsProject),
			},
		},
		{
			Name: "install available app by name",
			Test: c[1].installAvailableApp(jsProject),
			Require: test.Tests{
				c[0].userClaim(c[1]),
				c[0].publishProject(jsProject),
			},
		},
		{
			Name:    "run js app",
			Test:    c[1].runApp(jsProject),
			Require: test.Tests{c[1].installAvailableApp(jsProject)},
		},
		{
			Name: "create js-rollup project",
			Test: c[0].newProject(jsRollupProject),
		},
		{
			Name: "build js-rollup project",
			Test: c[0].buildProject(jsRollupProject),
		},
		{
			Name: "publish js-rollup app",
			Test: c[0].publishProject(jsRollupProject),
		},
		{
			Name: "install js-rollup app by name",
			Test: c[1].installAvailableApp(jsRollupProject),
			Require: test.Tests{
				c[0].userClaim(c[1]),
				c[0].publishProject(jsRollupProject),
			},
		},
		{
			Name:    "run js-rollup app",
			Test:    c[1].runApp(jsRollupProject),
			Require: test.Tests{c[1].installAvailableApp(jsRollupProject)},
		},
		// ====== dev.html ======
		{
			Name: "list html templates",
			Test: c[0].listTemplates("dev.html"),
		},
		{
			Name: "create html project",
			Test: c[0].newProject(htmlProject),
		},
		{
			Name: "create svelte project",
			Test: c[0].newProject(svelteProject),
		},
		{
			Name: "build svelte project",
			Test: c[0].buildProject(svelteProject),
		},
		{
			Name: "create react project",
			Test: c[0].newProject(reactProject),
		},
		{
			Name: "build react project",
			Test: c[0].buildProject(reactProject),
		},
		// ====== tear down ======
		{
			Name:  "portal close",
			Test:  c[0].portalClose(),
			Group: 1,
		},
		{
			Name:  "print logs",
			Group: 2,
			Test:  c[0].printLog(),
		},
		{
			Name:  "print logs",
			Group: 3,
			Test:  c[1].printLog(),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d  %s", i, tt.Name), runner.Run(tests, tt))
	}

	t.Cleanup(func() {
		time.Sleep(1 * time.Second) // await logs
		forceStopContainers(c...)
	})
}

var jsProject = projectOpts{
	runner:   "dev.js",
	template: "js",
	name:     "js-raw",
}

var jsRollupProject = projectOpts{
	runner:   "dev.js",
	template: "js-rollup",
}

var htmlProject = projectOpts{
	runner:   "dev.html",
	template: "html",
}

var svelteProject = projectOpts{
	runner:   "dev.html",
	template: "svelte",
}

var reactProject = projectOpts{
	runner:   "dev.html",
	template: "react",
}
