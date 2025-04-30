package apps

import (
	"github.com/cryptopunkscc/portal/test"
	"testing"
)

func TestAppsBuild_Run(t *testing.T) {
	err := Build("clean", "pack")
	test.AssertErr(t, err)
}
