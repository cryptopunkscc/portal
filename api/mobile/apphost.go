package mobile

//type Apphost interface {
//	Resolve(name string) (s string, err error)
//	Register(name string) (ApphostListener, error)
//	Query(nodeID string, query string) (c Conn, err error)
//}
//
//type ApphostListener interface {
//	Next() (query QueryData, err error)
//}
//
//type QueryData interface {
//	Caller() string
//	Accept() (c Conn, err error)
//	Reject() error
//	Query() string
//}

type Conn interface {
	Reader
	Writer
	Closer
}

type ReadCloser interface {
	Reader
	Closer
}

type Reader interface {
	Read(p []byte) (n int, err error)
	ReadN(n int) (arr []byte, err error)
	ReadAll() ([]byte, error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Closer interface {
	Close() error
}
