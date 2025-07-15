package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

type output struct {
	out  io.Writer
	buf  bytes.Buffer
	errs []error
}

var errorPattern = "ERROR: "

func (o *output) Write(p []byte) (int, error) {
	for start := 0; start < len(p); {

		idx := bytes.IndexByte(p[start:], '\n')

		if idx < 0 {
			o.buf.Write(p[start:])
			break
		}

		end := start + idx + 1
		o.buf.Write(p[start:end])
		line := o.buf.String()

		if strings.HasPrefix(line, errorPattern) {
			err := errors.New(line[len(errorPattern):])
			o.errs = append(o.errs, err)
		}

		o.buf.Reset()
		start = end
	}

	return o.out.Write(p)
}

func (o *output) Error() error {
	return errors.Join(o.errs...)
}

func (o *output) Exit() {
	os.Exit(len(o.errs))
}
