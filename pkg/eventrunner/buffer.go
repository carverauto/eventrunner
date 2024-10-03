// Package eventrunner pkg/eventrunner/buffer.go
package eventrunner

import "bytes"

func (er *EventRouter) getBuffer() *bytes.Buffer {
	return er.bufferPool.Get().(*bytes.Buffer)
}

func (er *EventRouter) putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	er.bufferPool.Put(buf)
}
