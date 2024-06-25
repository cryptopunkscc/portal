package notify

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/android"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

type testClient struct {
	Client
}

func NewTestClient() ApiClient {
	return &Client{port: testPort}
}

func ConnectTestClient(t *testing.T) (ApiClient, func()) {
	c := &Client{port: testPort}
	if err := c.Connect(); err != nil {
		t.Fatal(err)
	}
	c.Logger(log.Default())
	return c, func() {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestServer(returnErr bool) (cancelFunc context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	if err := rpc.NewApp(testPort).Interface(testService{err: returnErr}).Run(ctx); err != nil {
		return
	}
	time.Sleep(time.Second * 1)
	return
}

const testPort = "android/notify/jrpc/test"

var _ android.NotifyServiceApi = testService{}

type testService struct {
	err bool
}

func (t testService) String() string {
	return testPort
}

func (t testService) Create(channel *android.Channel) (err error) {
	if t.err {
		err = Response(channel)
	} else {
		log.Println(channel)
	}
	return
}

func (t testService) Notify(notification *android.Notification) (err error) {
	if t.err {
		err = Response(notification)
	} else {
		log.Println(notification)
	}
	return
}

func Response(args ...any) error {
	return fmt.Errorf("%v", args...)
}

func Verify(t *testing.T, err error, args ...any) {
	assert.EqualError(t, err, fmt.Sprint(args...))
}
