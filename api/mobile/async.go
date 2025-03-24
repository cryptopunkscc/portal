package mobile

type AsyncReader struct {
	reader Reader
}

func NewAsyncReader(reader Reader) *AsyncReader {
	return &AsyncReader{reader: reader}
}

func (c *AsyncReader) Read(p []byte, callback CallbackN) {
	go func() { callback.Result(c.reader.Read(p)) }()
}

func (c *AsyncReader) ReadAll(callback CallbackBytes) {
	go func() { callback.Result(c.reader.ReadAll()) }()
}

type AsyncWriter struct {
	writer Writer
}

func NewAsyncWriter(writer Writer) *AsyncWriter {
	return &AsyncWriter{writer: writer}
}

func (c *AsyncWriter) Write(p []byte, callback CallbackN) {
	go func() { callback.Result(c.writer.Write(p)) }()
}

type CallbackN interface {
	Result(n int, err error)
}

type CallbackBytes interface {
	Result(arr []byte, err error)
}
