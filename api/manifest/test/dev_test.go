package test

import (
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDev_UnmarshalFrom(t *testing.T) {
	expected := manifest.Dev{
		Dist: manifest.Dist{
			App:     manifest.App{Name: "name", Description: "description", Version: 0},
			Api:     manifest.Api{Version: 1},
			Config:  manifest.Config{Timeout: 0, Hidden: false},
			Release: manifest.Release{Version: 0},
		},
		Builds: manifest.Builds{Builds: manifest.InnerBuilds{
			Build: manifest.Build{Cmd: "cmd2", Exec: "exec2"},
			Builds: map[string]manifest.InnerBuilds{
				"os1": {},
				"os2": {
					Build: manifest.Build{Out: "main3", Exec: "exec3"},
					Builds: map[string]manifest.InnerBuilds{
						"arch1": {},
						"arch2": {Build: manifest.Build{Out: "main4", Cmd: "cmd4"}},
					},
				},
			},
		}},
	}
	actual := manifest.Dev{}
	err := actual.UnmarshalFrom(DevYml)
	test.AssertErr(t, err)
	assert.Equal(t, expected, actual)
}
