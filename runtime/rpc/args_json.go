package rpc

import (
	"encoding/json"
	"io"
)

type jsonArgsDecoder struct{}

func NewJsonArgsDecoder() ArgsDecoder {
	return &jsonArgsDecoder{}
}

func (d jsonArgsDecoder) Test(b []byte) bool {
	return b[0] == '[' || b[0] == '{'
}

func (d jsonArgsDecoder) TestScan(scan io.ByteScanner) bool {
	b, err := scan.ReadByte()
	_ = scan.UnreadByte()
	return err == nil && (b == '[' || b == '{')
}

func (d jsonArgsDecoder) Unmarshal(bytes []byte, args []any) error {
	if len(args) == 1 {
		// unmarshal struct payload to as first arg
		if bytes[0] == '{' {
			return json.Unmarshal(bytes, &args[0])
		}
	}

	// perform default unmarshal
	return json.Unmarshal(bytes, &args)
}

func (d jsonArgsDecoder) Decode(conn ByteScannerReader, args []any) error {
	jd := json.NewDecoder(conn)
	if len(args) == 1 {
		// unmarshal struct payload to as first arg
		c, err := conn.ReadByte()
		_ = conn.UnreadByte()
		if err != nil {
			return err
		}
		if c == byte('{') {
			return jd.Decode(&args[0])
		}
	}

	// perform default unmarshal
	return jd.Decode(&args)
}
