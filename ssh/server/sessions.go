// Copyright 2021-2023 arcadium.dev <info@arcadium.dev>
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

package server // import "arcadium.dev/core/ssh/server"

import (
	"fmt"
	"sync"
)

type (
	// Sessions manages the sessions. It provides a way to coordinate the
	// closing of all sessions.
	Sessions struct {
		mu       sync.RWMutex
		wg       sync.WaitGroup
		closing  bool
		sessions map[*Session]struct{}
	}
)

// NewSession returns a new Sessions object.
func NewSessions() *Sessions {
	return &Sessions{
		sessions: make(map[*Session]struct{}),
	}
}

func Type() string {
	return "session"
}

// Add adds the session. If the sessions are in the process of closing and
// error will be returned.
func (s *Sessions) Add(session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closing {
		return fmt.Errorf("cannot add session, closing")
	}

	s.sessions[session] = struct{}{}
	s.wg.Add(1)

	return nil
}

// Remove removes the given session. If the sessions are in the process of
// closing an error will be returned.
func (s *Sessions) Remove(session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closing {
		return fmt.Errorf("cannot remove session, closing")
	}

	delete(s.sessions, session)
	s.wg.Done()

	return nil
}

// Close closes all sessions.
func (s *Sessions) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.closing = true

	for session := range s.sessions {
		session.Close()
		delete(s.sessions, session)
		s.wg.Done()
	}
}

// Wait blocks until all sessions have been closed.
func (s *Sessions) Wait() {
	s.wg.Wait()
}
