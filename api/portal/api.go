package portal

type Client interface {
	Join()
	Ping() (err error)
	Open(src ...string) error
	Close() error
}
