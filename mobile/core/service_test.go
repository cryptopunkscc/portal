package core

import (
	"errors"
	"io/fs"
	"net"
	"testing"
	"time"

	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/query"
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	ip "github.com/cryptopunkscc/astrald/mod/ip/src"
	"github.com/cryptopunkscc/portal/mobile"
	"github.com/cryptopunkscc/portal/pkg/util/test"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	api := &testApi{T: t, dir: test.CleanMkdir(t), event: make(chan any, 100)}
	srv := Create(api).(*Service)
	ip.InterfaceAddrs = net.InterfaceAddrs
	ether.NetInterfaces = ether.DefaultNetInterfaces

	srv.Start()
	api.requireEvent(mobile.STARTING, 2*time.Second)
	api.requireEvent(mobile.STARTED, 5*time.Second)

	ch, err := srv.client.QueryChannel(srv.ctx, "portal.open", query.Args{"app": "test_app"})
	require.NoError(t, err)
	err = ch.Switch(channel.ExpectAck, channel.PassErrors, channel.WithContext(srv.ctx))
	require.Equal(t, err.Error(), fs.ErrNotExist.Error())

	srv.Stop()
	api.requireEvent("context canceled", 2*time.Second)
	api.requireEvent(mobile.STOPPED, 2*time.Second)
}

type testApi struct {
	*testing.T
	dir   string
	event chan any
}

func (t *testApi) requireEvent(event any, timeout time.Duration) {
	select {
	case e := <-t.event:
		require.Equal(t, event, e)
	case <-time.After(timeout):
		require.FailNow(t.T, "timeout")
	}
}

func (t *testApi) CacheDir() string {
	return t.dir
}

func (t *testApi) DataDir() string {
	return t.dir
}

func (t *testApi) DbDir() string {
	return t.dir
}

func (t *testApi) Status(id int32) {
	t.Logf("status: %d", id)
	t.event <- id
}

func (t *testApi) Error(message string) {
	t.Logf("error: %s", message)
	t.event <- message
}

func (t *testApi) StartHtml(pkg string, args string) error {
	t.Logf("startHtml: %s %s", pkg, args)
	t.event <- []string{pkg, args}
	return nil
}

func (t *testApi) Net() mobile.Net {
	return &testNet{}
}

type testNet struct{}

func (t testNet) Addresses() (string, error) {
	return "", errors.New("not implemented")
}

func (t testNet) Interfaces() (mobile.NetInterfaceIterator, error) {
	return nil, errors.New("not implemented")
}
