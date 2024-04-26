package apphost

import "github.com/cryptopunkscc/astrald/lib/astral"

type Flat interface {
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
	NodeInfo(identity string) (info NodeInfo, err error)
}

func NewFlatAdapter() Flat {
	return &FlatAdapter{
		listeners:   map[string]*astral.Listener{},
		connections: map[string]*Conn{},
	}
}
