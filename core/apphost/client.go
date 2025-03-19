package apphost

import (
	"github.com/cryptopunkscc/astrald/lib/apphost"
	"os"
	"sync"
)

type Client struct {
	apphost.Client
	mu sync.Mutex
}

func (c *Client) IsConnected() bool {
	return c.HostID != nil
}

func (c *Client) Connect() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.IsConnected() {
		return
	}
	return c.connect()
}

func (c *Client) Reconnect() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connect()
}

func (c *Client) connect() (err error) {
	if len(c.Endpoint) == 0 {
		c.Endpoint = apphost.DefaultEndpoint
	}
	if len(c.AuthToken) == 0 {
		c.AuthToken = os.Getenv(apphost.AuthTokenEnv)
	}
	client, err := apphost.NewClient(c.Endpoint, c.AuthToken)
	if err == nil {
		c.Client = *client
	}
	return
}
