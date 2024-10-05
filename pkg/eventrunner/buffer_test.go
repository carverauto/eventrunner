package eventrunner

import (
	"bytes"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventRouter_getBuffer(t *testing.T) {
	er := &EventRouter{
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
	}

	buf := er.getBuffer()
	assert.NotNil(t, buf)
	assert.IsType(t, &bytes.Buffer{}, buf)
}

func TestEventRouter_putBuffer(t *testing.T) {
	er := &EventRouter{
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
	}

	buf := &bytes.Buffer{}
	er.putBuffer(buf)

	// Check that the buffer is reset
	assert.Equal(t, 0, buf.Len())
}
