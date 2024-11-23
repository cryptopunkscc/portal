package rpc

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
)

type Serializer struct {
	io.Reader
	io.Writer
	io.Closer
	Marshal
	Unmarshal

	scanner *bufio.Scanner
	reader  io.Reader
}

type Marshal func(v any) ([]byte, error)
type Unmarshal func(data []byte, v any) error
type Read func(p []byte) (n int, err error)
type Write func(p []byte) (n int, err error)
type Closer func() error

func (r Read) Read(p []byte) (n int, err error)   { return r(p) }
func (w Write) Write(p []byte) (n int, err error) { return w(p) }
func (c Closer) Close() error                     { return c() }

type Failure struct {
	Error string `json:"error"`
}

func NewSerializer(rwc io.ReadWriteCloser) *Serializer {
	return &Serializer{Reader: rwc, Writer: rwc, Closer: rwc}
}

func (s *Serializer) init() {
	if s.Reader == nil {
		s.Reader = Read(func(p []byte) (n int, err error) { return 0, io.EOF })
	}
	if s.Writer == nil {
		s.Writer = Write(func(p []byte) (n int, err error) { return len(p), nil })
	}
	if s.Closer == nil {
		s.Closer = Closer(func() error { return nil })
	}
	if s.reader != s.Reader {
		s.reader = s.Reader
		s.scanner = bufio.NewScanner(s.Reader)
	}
}

func (s *Serializer) Encode(value any) (err error) {
	s.init()
	r := value
	if value != End {
		switch v := value.(type) {
		case error:
			r = Failure{v.Error()}
		}
	}
	var data []byte
	data, err = s.Marshal(r)
	if err != nil || data == nil {
		if !errors.Is(err, End) {
			return
		}
	}
	data = append(data, '\n')
	_, err = s.Write(data)
	return
}

func (s *Serializer) Decode(value any) (err error) {
	s.init()

	bytes, err := s.Bytes()
	if err != nil {
		return
	}
	if len(bytes) == 0 {
		return End
	}

	// try decode as failure
	f := Failure{}
	if err = s.Unmarshal(bytes, &f); err == nil && f.Error != "" {
		return fmt.Errorf("RPC: %s", f.Error)
	}

	// decode value
	return s.Unmarshal(bytes, value)
}

func (s *Serializer) Bytes() (bytes []byte, err error) {
	s.init()
	if !s.scanner.Scan() {
		if err = s.scanner.Err(); err == nil {
			err = io.EOF
		}
		return
	}
	s.scanner.Bytes()
	bytes = s.scanner.Bytes()
	return
}

func (s *Serializer) Logger(logger plog.Logger) {

}
