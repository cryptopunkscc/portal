package rpc

import (
	"errors"
	"fmt"
	"github.com/leaanthony/clir"
	"io"
	"reflect"
	"strings"
)

type clirArgsDecoder struct{}

func NewClirArgsDecoder() ArgsDecoder { return &clirArgsDecoder{} }

func (d clirArgsDecoder) Decode(conn ByteScannerReader, args []any) (err error) {
	end := []byte("\n")
	b := byte(0)
	var s = 0
	var n = 0
	for {
		if b, err = conn.ReadByte(); err != nil {
			return
		}
		n++
		if b != end[s] {
			s = 0
			continue
		}
		s++
		if s == len(end) {
			break
		}
	}
	for i := 0; i < n; i++ {
		_ = conn.UnreadByte()
	}
	bytes := make([]byte, n)
	if _, err = conn.Read(bytes); err != nil {
		return
	}
	bytes = bytes[:len(bytes)-1]
	err = d.Unmarshal(bytes, args)
	return
}

func (d clirArgsDecoder) Test(b []byte) bool {
	return b[0] == '$' || b[0] == ' '
}

func (d clirArgsDecoder) TestScan(scan io.ByteScanner) bool {
	b, err := scan.ReadByte()
	_ = scan.UnreadByte()
	return err == nil && (b == '$' || b == ' ')
}

func (d clirArgsDecoder) Unmarshal(bytes []byte, args []any) (err error) {
	f := strings.Fields(string(bytes[1:]))
	for i, s := range f {
		if len(s) < 2 {
			continue
		}
		v := s[0]
		if v == '"' || v == '\'' {
			s = s[1:]
		}
		v = s[len(s)-1]
		if v == '"' || v == '\'' {
			s = s[:len(s)-1]
		}
		f[i] = s
	}
	c := clir.NewCli("", "", "").Action(func() error { return nil })
	flags := false
	for _, a := range args {
		v := reflect.ValueOf(a)
		if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
			c.AddFlags(a)
			flags = true
			continue
		}

		if len(f) == 0 {
			break
		}
		s := f[0]
		f = f[1:]
		_, err = fmt.Sscan(s, a)
		if err == nil {
			continue
		}

		return errors.New("invalid arg type")
	}
	if flags {
		err = c.Run(f...)
	}
	return err
}
