package stream

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"strings"
)

type Serializer struct {
	Codec

	io.Reader
	io.Writer
	io.Closer

	scanner *bufio.Reader
	reader  io.Reader
	log     plog.Logger
}

type Codec struct {
	MarshalArgs Marshal
	Marshal     Marshal
	Unmarshal   Unmarshal
	Ending      []byte
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
		s.scanner = bufio.NewReader(s.Reader)
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
	data = append(data, s.Ending...)
	_, err = s.Write(data)
	return
}

func (s *Serializer) Decode(value any) (err error) {
	s.init()

	b, err := s.ReadBytes('\n')
	if err != nil {
		return
	}
	if len(b) == 0 {
		return End
	}

	// try decode as failure
	f := Failure{}
	if err = s.Unmarshal(b, &f); err == nil && f.Error != "" {
		return fmt.Errorf("RPC: %s", f.Error)
	}

	// decode value
	return s.Unmarshal(b, value)
}

func (s *Serializer) ReadBytes(delim byte) (b []byte, err error) {
	s.init()
	if b, err = s.scanner.ReadBytes(delim); err != nil {
		return
	}
	b = bytes.TrimSuffix(b, []byte{delim})
	if s.log != nil {
		s.log.Printf("< %s [%db]", strings.Trim(string(b), "\n"), len(b))
	}
	return
}

func (s *Serializer) ReadString(delim byte) (str string, err error) {
	s.init()
	if str, err = s.scanner.ReadString(delim); err != nil {
		return
	}
	str = strings.TrimSuffix(str, string([]byte{delim}))
	if s.log != nil {
		s.log.Printf("< %s [%db]", strings.Trim(str, "\n"), len(str))
	}
	return
}

func (s *Serializer) Read(b []byte) (n int, err error) {
	n, err = s.Reader.Read(b)
	if s.log != nil {
		s.log.Printf("< %s [%db]", strings.Trim(string(b[:n]), "\n"), n)
	}
	return
}

func (s *Serializer) Write(b []byte) (n int, err error) {
	n, err = s.Writer.Write(b)
	if s.log != nil {
		s.log.Printf("> %s [%db]", strings.Trim(string(b[:n]), "\n"), n)
	}
	return
}

func (s *Serializer) Logger(log plog.Logger) {
	s.log = log
}
