package mobile

type Apphost interface {
	Resolve(name string) (s string, err error)
	Register(name string) (ApphostListener, error)
	Query(nodeID string, query string) (c Conn, err error)
}

type ApphostListener interface {
	Next() (query QueryData, err error)
}

type QueryData interface {
	Caller() string
	Accept() (c Conn, err error)
	Reject() error
	Query() string
}

type Conn interface {
	Read(p []byte) (n int, err error)
	ReadN(n int) (arr []byte, err error)
	Write(p []byte) (n int, err error)
	Close() error
}
