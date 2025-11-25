package test

import (
	"testing"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestDist_UnmarshalFrom(t *testing.T) {
	expected := manifest.Dist{
		App: manifest.App{
			Name:        "name1",
			Description: "description1",
		},
		Api: manifest.Api{
			Version: 1,
			Handler: cmd.Handler{
				Desc: "desc2",
				Params: []cmd.Param{
					{Type: "string"},
					{Name: "name", Type: "string", Desc: "desc"},
				},
			},
		},
	}
	actual := manifest.Dist{}
	err := actual.UnmarshalFrom(DistYml)
	test.AssertErr(t, err)
	assert.Equal(t, expected, actual)
}
