// Copyright 2021-2022 arcadium.dev <info@arcadium.dev>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log // import "arcadium.dev/core/log

import (
	"sync"
)

type (
	// StringBuffer implements a simple buffer that can be used in tests.  It
	// implements the io.Writer interface. Each write will append a string
	// to the Buffer.
	StringBuffer struct {
		lock   sync.RWMutex
		buffer []string
	}
)

// NewStringBuffer returns a StringBuffer.
func NewStringBuffer() *StringBuffer {
	return &StringBuffer{buffer: make([]string, 0)}
}

// Write implements the io.Writer interface.
func (l *StringBuffer) Write(p []byte) (int, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.buffer = append(l.buffer, string(p))
	return len(p), nil
}

func (l *StringBuffer) Len() int {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return len(l.buffer)
}

func (l *StringBuffer) Index(i int) string {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.buffer[i]
}
