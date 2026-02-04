package mobile

type Assets interface {
	Get(uri string) (Asset, error)
}

type Asset interface {
	Mime() string
	Encoding() string
	Data() ReadCloser
}
