package win10

import (
	"fmt"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
)

func Test_PrepareWindowsVM(t *testing.T) {
	runner := test.Runner{}
	tests := []test.Task{
		{
			Name: "remove VM",
			Test: removeVmDisk(),
			Require: test.Tests{
				virtUndefine(),
			},
		},
		{
			Name: "create windows virtual machine",
			Test: virtInstall(),
			Require: test.Tests{
				genSshKey(),
				copyPubKey(),
				genCiDataIso(),
				genAutounattendIso(),
				removeVmDisk(),
				createQemuImage(),
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d  %s", i, tt.Name), runner.Run(tests, tt))
	}
}
