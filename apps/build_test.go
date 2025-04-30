package apps

import (
	"github.com/cryptopunkscc/portal/runner/apps_build"
	"github.com/cryptopunkscc/portal/test"
	"testing"
)

func TestAppsBuild_Run(t *testing.T) {
	err := apps_build.Run("clean", "pack")
	test.AssertErr(t, err)
}
