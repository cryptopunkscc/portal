package rpc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
	"github.com/cryptopunkscc/portal/pkg/util/rpc/caller/query"
	"github.com/cryptopunkscc/portal/pkg/util/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/util/rpc/cmd"
	"github.com/cryptopunkscc/portal/pkg/util/rpc/router"
	"github.com/stretchr/testify/assert"
)

func TestRouter_routeQuery(t *testing.T) {
	type fields struct {
		Base   router.Base
		Logger plog.Logger
	}
	type args struct {
		q *testPendingQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		setup   func(t *testing.T, f *fields, a *args)
		verify  func(t *testing.T, f fields, a args)
	}{
		{
			name:    "should call function and return result",
			wantErr: io.EOF,
			fields: fields{
				Logger: plog.New(),
				Base: router.Base{
					Registry: router.CreateRegistry(cmd.Handler{
						Sub: cmd.Handlers{cmdHandler},
					}),
					Unmarshal: query.Unmarshal,
				},
			},
			args: args{
				q: &testPendingQuery{
					testConn: &testConn{
						query:       "cmd?_=1",
						ReadBuffer:  bytes.NewBuffer(nil),
						WriteBuffer: bytes.NewBuffer(nil),
					},
				},
			},
			setup: func(t *testing.T, f *fields, a *args) {
				a.q.testConn.Writer = bufio.NewWriter(a.q.testConn.WriteBuffer)
				a.q.testConn.Reader = bufio.NewReader(a.q.testConn.ReadBuffer)
			},
			verify: func(t *testing.T, f fields, a args) {
				assert.Equal(t, "{\"value\":1}\n", a.q.testConn.WriteBuffer.String())
			},
		},
		{
			name:    "should change encoder",
			wantErr: io.EOF,
			fields: fields{
				Logger: plog.New(),
				Base: router.Base{
					Registry: router.CreateRegistry(cmd.Handler{
						Sub: cmd.Handlers{
							cli.EncodingHandler,
							cmdHandler,
						},
					}),
					Unmarshal: query.Unmarshal,
				},
			},
			args: args{
				q: &testPendingQuery{
					testConn: &testConn{
						query:       "encoding?format=cli&encoder=cli",
						ReadBuffer:  bytes.NewBufferString("cmd 1\n"),
						WriteBuffer: bytes.NewBuffer(nil),
					},
				},
			},
			setup: func(t *testing.T, f *fields, a *args) {
				a.q.testConn.Writer = bufio.NewWriter(a.q.testConn.WriteBuffer)
				a.q.testConn.Reader = bufio.NewReader(a.q.testConn.ReadBuffer)
			},
			verify: func(t *testing.T, f fields, a args) {
				assert.Equal(t, "1\n", a.q.testConn.WriteBuffer.String())
			},
		},
		{
			name: "should switch to cli",
			fields: fields{
				Logger: plog.New(),
				Base: router.Base{
					Registry: router.CreateRegistry(cmd.Handler{
						Sub: cmd.Handlers{
							cli.Handler,
							cmdHandler,
						},
					}),
					Unmarshal: query.Unmarshal,
				},
			},
			args: args{
				q: &testPendingQuery{
					testConn: &testConn{
						query:       "cli",
						ReadBuffer:  bytes.NewBufferString("cmd 1\n"),
						WriteBuffer: bytes.NewBuffer(nil),
					},
				},
			},
			setup: func(t *testing.T, f *fields, a *args) {
				a.q.testConn.Writer = bufio.NewWriter(a.q.testConn.WriteBuffer)
				a.q.testConn.Reader = bufio.NewReader(a.q.testConn.ReadBuffer)
			},
			verify: func(t *testing.T, f fields, a args) {
				assert.Equal(t, "1\n", a.q.testConn.WriteBuffer.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}

			r := &Router{}
			r.Base = tt.fields.Base
			r.Log = tt.fields.Logger

			if err := r.routeQuery(tt.args.q); !errors.Is(err, tt.wantErr) {
				t.Errorf("routeQuery() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.verify != nil {
				tt.verify(t, tt.fields, tt.args)
			}
		})
	}
}

var _ PendingQuery = &testPendingQuery{}

type testPendingQuery struct {
	identity  *astral.Identity
	rejectErr error
	testConn  *testConn
	conn      *apphost.Conn
	acceptErr error
	closeErr  error
}

func (t *testPendingQuery) Query() string                    { return t.conn.Query().Query }
func (t *testPendingQuery) RemoteIdentity() *astral.Identity { return t.identity }
func (t *testPendingQuery) Reject() error                    { return t.rejectErr }
func (t *testPendingQuery) Accept() *apphost.Conn            { return t.conn }
func (t *testPendingQuery) Close() error                     { return t.closeErr }

var _ net.Conn = &testConn{}

type testConn struct {
	*bufio.Writer
	*bufio.Reader
	ReadBuffer  *bytes.Buffer
	WriteBuffer *bytes.Buffer
	closeErr    error
	remoteAddr  net.Addr
	query       string
}

func (c *testConn) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (c *testConn) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *testConn) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *testConn) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *testConn) Write(b []byte) (int, error) {
	defer c.Writer.Flush()
	return c.Writer.Write(b)
}

func (c *testConn) Close() error { return c.closeErr }

// func (c *testConn) RemoteIdentity() *astral.Identity { return c.RemoteIdentity() }
func (c *testConn) RemoteAddr() net.Addr { return c.remoteAddr }

//func (c *testConn) Query() string                    { return c.query }

var cmdHandler = cmd.Handler{
	Name: "cmd",
	Func: func(v int) cmdResponse { return cmdResponse{Value: v} },
}

type cmdResponse struct {
	Value int `json:"value"`
}

func (c cmdResponse) MarshalCLI() string {
	return fmt.Sprintf("%d", c.Value)
}
