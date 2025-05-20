package apps

import (
	"github.com/cryptopunkscc/portal/pkg/test"
	"testing"
)

func TestAppsBuild_Run(t *testing.T) {
	err := Build("clean", "pack")
	if err != nil && err.Error() != "npm is required but not installed" {
		test.AssertErr(t, err)
	}
}
