package rpc

import (
	"errors"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"reflect"
)

type Serializer struct {
	io.WriteCloser
	ByteScannerReader
	logger    *ConnLogger
	enc       Encoder
	dec       Decoder
	marshal   Marshal
	unmarshal Unmarshal
	remoteID  id.Identity
	codecs    Codecs
}

type Encoder interface{ Encode(v any) error }
type Decoder interface{ Decode(v any) error }
type Marshal func(v any) ([]byte, error)
type Unmarshal func(data []byte, v any) error
type RemoteIdInfo interface{ RemoteIdentity() id.Identity }
type Codecs func(io.ReadWriter) (Encoder, Decoder, Marshal, Unmarshal)
type raw struct{ bytes []byte }

func (s *Serializer) RemoteIdentity() (i id.Identity) {
	return s.remoteID
}

func (s *Serializer) Codecs(codecs Codecs) {
	s.codecs = codecs
}

func (s *Serializer) Logger(logger plog.Logger) {
	s.setLogger(logger)
	s.setupEncoding()
}

func (s *Serializer) setLogger(logger plog.Logger) {
	if s.logger == nil {
		s.logger = NewConnLogger(s, logger)
	} else {
		s.logger.Logger = logger
	}
}

func (s *Serializer) setupEncoding() {
	if s.codecs == nil {
		s.codecs = JsonCodecs
	}
	var rw io.ReadWriter = s
	if s.logger != nil {
		rw = s.logger
	}
	s.enc, s.dec, s.marshal, s.unmarshal = s.codecs(rw)
}

func (s *Serializer) setConn(conn io.ReadWriteCloser) {
	s.WriteCloser = conn
	s.ByteScannerReader = NewByteScannerReader(conn)
	if s.logger != nil {
		s.logger.ReadWriteCloser = s
	}
	s.setupEncoding()
	s.setupRemoteID()
}

func (s *Serializer) setupRemoteID() {
	var wc any = s.WriteCloser
	if info, ok := wc.(RemoteIdInfo); ok && !reflect.ValueOf(info).IsNil() {
		s.remoteID = info.RemoteIdentity()
	}
}

func (s *Serializer) Encode(value any) (err error) {
	r := value
	switch v := value.(type) {
	case error:
		r = Failure{v.Error()}
	}
	return s.enc.Encode(r)
}

func (s *Serializer) Decode(value any) (err error) {
	// decode raw value
	r := raw{}
	if err = s.dec.Decode(&r); err != nil {
		return
	}

	// try decode as failure
	f := Failure{}
	if err = s.unmarshal(r.bytes, &f); err == nil && f.Error != "" {
		return errors.New(f.Error)
	}

	// decode value
	return s.unmarshal(r.bytes, value)
}
