package rpc

import (
	"errors"
	"io"
)

type ArgsDecoder interface {
	Test([]byte) bool
	TestScan(scanner io.ByteScanner) bool
	Unmarshal(bytes []byte, args []any) error
	Decode(conn ByteScannerReader, args []any) error
}

type argsDecoders struct {
	decoders []ArgsDecoder
}

func (a *argsDecoders) Append(decoders []ArgsDecoder) *argsDecoders {
	a.decoders = append(a.decoders, decoders...)
	return a
}

func (a *argsDecoders) TestScan(scan io.ByteScanner) (b bool) {
	return a.find(scan) != nil
}

func (a *argsDecoders) find(scan io.ByteScanner) (d ArgsDecoder) {
	for _, decoder := range a.decoders {
		if ok := decoder.TestScan(scan); ok {
			return decoder
		}
	}
	return
}

func (a *argsDecoders) Decode(scanner ByteScannerReader, args []any) (err error) {
	// find decoder
	d := a.find(scanner)

	// try decode
	if d != nil {
		err = d.Decode(scanner, args)
		return
	}

	// drop unknown bytes from scanner
	scanner.Clear()
	err = errors.New("cannot decode unknown format")
	return
}

func (a *argsDecoders) Test(bytes []byte) bool {
	for _, decoder := range a.decoders {
		if decoder.Test(bytes) {
			return true
		}
	}
	return false
}

func (a *argsDecoders) Unmarshal(bytes []byte, args []any) error {
	for _, decoder := range a.decoders {
		if decoder.Test(bytes) {
			return decoder.Unmarshal(bytes, args)
		}
	}
	return errors.New("cannot unmarshal unknown format")
}
