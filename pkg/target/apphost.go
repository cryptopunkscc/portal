package target

import "github.com/cryptopunkscc/go-astral-js/pkg/apphost"

type Apphost interface {
	Prefix() []string
	Close() error
	Interrupt()
	Log(arg ...any)
	LogArr(arg []any)
	Sleep(duration int64)
	ServiceRegister(service string) (err error)
	ServiceClose(service string) (err error)
	ConnAccept(service string) (data string, err error)
	ConnClose(id string) (err error)
	ConnWrite(id string, data string) (err error)
	ConnRead(id string) (data string, err error)
	Query(identity string, query string) (data string, err error)
	QueryName(name string, query string) (data string, err error)
	Resolve(name string) (id string, err error)
	NodeInfo(identity string) (info apphost.NodeInfo, err error)
}
