package test

import (
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilds(t *testing.T) {
	b := &manifest.Builds{}
	err := b.UnmarshalFrom(BuildYml)
	test.AssertErr(t, err)

	t.Run("targets", func(t *testing.T) {
		expected := [][]string{{"os1"}, {"os2", "arch1"}, {"os2", "arch2"}}
		actual := b.Targets()
		assert.ElementsMatch(t, expected, actual)
	})
	t.Run("get", func(t *testing.T) {
		t.Run("default", func(t *testing.T) {
			expectedBuild := manifest.Build{
				Out: "main1", Cmd: "cmd2", Exec: "exec2",
				Target: manifest.Target{Exec: "main1", OS: "linux", Arch: "amd64"},
			}
			actualBuild := b.Get()
			assert.Equal(t, expectedBuild, actualBuild)
		})
		t.Run("custom", func(t *testing.T) {
			expectedBuild := manifest.Build{
				Out: "main4", Cmd: "cmd4", Exec: "exec3",
				Target: manifest.Target{Exec: "main4", OS: "os2", Arch: "arch2"},
			}
			actualBuild := b.Get("os2", "arch2")
			assert.Equal(t, expectedBuild, actualBuild)
		})
	})
}
