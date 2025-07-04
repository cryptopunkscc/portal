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
		{
			Name: "start portald via portal",
			Test: c[0].portalStart(),
		},
		{
			Name: "help",
			Test: c[0].portalHelp(),
		},
		{
			Name: "create user",
			Test: c[0].createUser(),
		},
		{
			Name: "claim",
			Test: c[0].claim(c[1]),
		},
		{
			Name: "user info",
			Test: c[0].userInfo(),
		},
		{
			Name:  "close",
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
