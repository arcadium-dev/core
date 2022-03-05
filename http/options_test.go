// Copyright 2021 arcadium.dev <info@arcadium.dev>
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

package http // import "arcadium.dev/core/http"

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"arcadium.dev/core/log"
)

func TestWithServerAddr(t *testing.T) {
	addr := "example.com:4201"
	s := &Server{}
	WithServerAddr(addr).apply(s)

	if s.addr != addr {
		t.Error("failed to set server addr")
	}
}

func TestWithServerTLS(t *testing.T) {
	s := &Server{
		server: &http.Server{},
	}
	cfg := &tls.Config{}
	WithServerTLS(cfg).apply(s)

	if s.server.TLSConfig == nil {
		t.Error("failed to set TLSConfig")
	}
}

func TestWithServerShutdownTimeout(t *testing.T) {
	s := &Server{}
	timeout := 52 * time.Second
	WithServerShutdownTimeout(timeout).apply(s)

	if s.shutdownTimeout != timeout {
		t.Errorf("Unexpected timeout: %v", s.shutdownTimeout)
	}
}

func TestWithServerLogger(t *testing.T) {
	s := &Server{}
	logger, err := log.New(log.WithLevel(log.LevelDebug), log.WithFormat(log.FormatLogfmt))
	if err != nil {
		t.Fatal("Failed to create logger")
	}
	WithServerLogger(logger).apply(s)

	if s.logger != logger {
		t.Errorf("Unexpected logger: %+v", logger)
	}
}
