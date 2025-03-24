package core

import (
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	tcp "github.com/cryptopunkscc/astrald/mod/tcp/src"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/test"
	"io"
	"log"
	"net"
	"testing"
)

func Test_runtime_Install(t *testing.T) {
	dir := test.Dir(t)

	core := Create(testApi{testDir: dir})

	tcp.InterfaceAddrs = net.InterfaceAddrs
	ether.NetInterfaces = ether.DefaultNetInterfaces

	if err := core.Install(); err != nil {
		plog.New().Println(err)
		t.Error(err)
	}

	app := core.App("portal.launcher")
	ass := app.Assets()
	indexAsset, err := ass.Get("index.html")
	if err != nil {
		plog.New().Println(err)
		t.Error(err)
		return
	}
	indexBytes, err := io.ReadAll(indexAsset.Data())
	if err != nil {
		plog.New().Println(err)
		t.Error(err)
		return
	}
	log.Println(string(indexBytes))
}
