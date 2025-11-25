package apps

import (
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
)

func TestAppsBuild_Run(t *testing.T) {
	err := Build("clean", "pack")
	if err != nil && err.Error() != "npm is required but not installed" {
		test.AssertErr(t, err)
	}
}
