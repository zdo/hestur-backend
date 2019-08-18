package hestur

import (
	"bytes"
	"encoding/binary"
)

type serverBuffer struct {
	buf *bytes.Buffer
}

// NewServerBuffer creates new server buffer.
func newServerBuffer() serverBuffer {
	serverBuffer := serverBuffer{}
	serverBuffer.buf = new(bytes.Buffer)
	return serverBuffer
}

func (b serverBuffer) Write(data interface{}) {
	binary.Write(b.buf, binary.BigEndian, data)
}

func (b serverBuffer) Bytes() []byte {
	return b.buf.Bytes()
}
