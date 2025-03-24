package core

import "io"

type reader struct{ io.Reader }

func (r reader) Read(arr []byte) (n int, err error) {
	n, err = r.Reader.Read(arr)
	return
}

func (r reader) ReadN(n int) (arr []byte, err error) {
	var l int
	arr = make([]byte, n)
	if l, err = r.Reader.Read(arr); err == nil {
		arr = arr[:l]
	}
	return
}

func (r reader) ReadAll() (all []byte, err error) {
	all, err = io.ReadAll(r)
	return all, err
}
