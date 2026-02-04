package js

import (
	"testing"

	"github.com/cryptopunkscc/portal/pkg/util/plog"
	"github.com/cryptopunkscc/portal/pkg/util/test"
)

func TestBuildPortalLib(t *testing.T) {
	err := BuildPortalLib()
	if err != nil && err.Error() == "npm is required but not installed" {
		plog.Println(err)
		t.Skip()
	} else {
		test.AssertErr(t, err)
	}
}
