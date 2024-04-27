package rpc

import (
	"io"
	"strings"
)

type ByteScannerReader interface {
	io.Reader
	io.ByteScanner
	Append(bytes []byte)
	Clear()
	IsEmpty() bool
	Buffer() []byte
}

type byteScannerReader struct {
	io.Reader
	offset int
	end    int
	buff   []byte
}

func (r *byteScannerReader) IsEmpty() bool {
	return r.end == 0
}

func (r *byteScannerReader) Clear() {
	r.offset = 0
	r.end = 0
}

func (r *byteScannerReader) Buffer() []byte {
	return r.buff
}

func NewByteScannerReader(reader io.Reader) ByteScannerReader {
	if reader == nil {
		reader = strings.NewReader("")
	}
	r := &byteScannerReader{Reader: reader}
	return r
}

func (r *byteScannerReader) ReadByte() (b byte, err error) {
	if r.offset == cap(r.buff) {
		l := cap(r.buff)*2 + 32
		r.buff = append(r.buff, make([]byte, l)...)
	}
	if r.offset == r.end {
		l := 0
		if l, err = r.Reader.Read(r.buff[r.end:cap(r.buff)]); err != nil {
			//if l, err = r.Reader.Read(r.buff); err != nil {
			return
		}
		r.end = r.offset + l
	}

	b = r.buff[r.offset]
	r.offset++
	return
}

func (r *byteScannerReader) Append(bytes []byte) {
	r.buff = append(r.buff[:r.end], bytes...)
	r.end = len(r.buff)
}

func (r *byteScannerReader) UnreadByte() error {
	if r.offset > 0 {
		r.offset--
	}
	return nil
}

func (r *byteScannerReader) Read(p []byte) (n int, err error) {
	if r.offset == r.end {
		r.offset = 0
		r.end = 0
		r.buff = nil
		return r.Reader.Read(p)
	}
	for n = 0; n < cap(p) && r.offset < r.end; n++ {
		p[n] = r.buff[r.offset]
		r.offset++
	}
	return
}

var _ io.ByteScanner = &byteScannerReader{}
