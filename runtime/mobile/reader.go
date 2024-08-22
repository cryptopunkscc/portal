package runtime

import "io"

type reader struct{ io.Reader }

func (r reader) ReadN(n int) (arr []byte, err error) {
	var l int
	arr = make([]byte, n)
	if l, err = r.Read(arr); err == nil {
		arr = arr[:l]
	}
	return
}

func (r reader) ReadAll() ([]byte, error) { return io.ReadAll(r) }
