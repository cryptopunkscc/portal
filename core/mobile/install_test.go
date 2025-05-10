package core

import (
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	tcp "github.com/cryptopunkscc/astrald/mod/tcp/src"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"io"
	"log"
	"net"
	"testing"
)

func Test_runtime_Install(t *testing.T) {
	t.SkipNow() //FIXME

	dir := test.Dir(t)

	core := Create(testApi{testDir: dir})

	tcp.InterfaceAddrs = net.InterfaceAddrs
	ether.NetInterfaces = ether.DefaultNetInterfaces

	if err := core.Install(); err != nil {
		plog.New().Println(err)
		t.Error(err)
	}

	tests := []struct {
		name string
	}{
		{name: "portal.launcher"},
		{name: "astrald.profile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := core.App(tt.name)
			ass := app.Assets()
			indexAsset, err := ass.Get("index.html")
			if err != nil {
				plog.New().Println(err)
				t.FailNow()
				return
			}
			indexBytes, err := io.ReadAll(indexAsset.Data())
			if err != nil {
				plog.New().Println(err)
				t.Error(err)
				return
			}
			log.Println(string(indexBytes))
		})
	}
}
