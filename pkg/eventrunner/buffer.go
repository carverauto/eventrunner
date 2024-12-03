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
