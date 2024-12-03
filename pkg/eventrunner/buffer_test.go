/*
* Copyright 2024 Carver Automation Corp.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
*  limitations under the License.
 */

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
