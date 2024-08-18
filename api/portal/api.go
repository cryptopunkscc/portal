package portal

type Client interface {
	Ping() (err error)
	Await()
	Open(src string) error
	Close() error
}
