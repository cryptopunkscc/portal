package android

const (
	ContentPort = "android/content/jrpc"
)

type Info struct {
	Uri  string
	Size int64
	Mime string
	Name string
}

type ContentServiceApi interface {
	String() string
	Info(uri string) (files *Info, err error)
	Reader(uri string, offset int64) (err error)
}
