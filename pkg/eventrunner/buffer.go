// Package eventrunner pkg/eventrunner/buffer.go
package eventrunner

import "bytes"

func (er *EventRouter) getBuffer() *bytes.Buffer {
	return er.bufferPool.Get().(*bytes.Buffer)
}

func (er *EventRouter) putBuffer(buf Buffer) {
	if bb, ok := buf.(*bytes.Buffer); ok {
		bb.Reset()
		er.bufferPool.Put(bb)
	}
}
