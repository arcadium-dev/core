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
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/mux"

	"arcadium.dev/core/config"
	"arcadium.dev/core/config/mock"
	"arcadium.dev/core/log"
	"arcadium.dev/core/test"
)

func TestServerNew(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		s := NewServer()

		if s.addr != defaultAddr {
			t.Errorf("\nExpected addr: %s\nActual addr:   %s", defaultAddr, s.addr)
		}
		if s.shutdownTimeout != defaultShutdownTimeout {
			t.Errorf("\nExpected timeout: %+v\nActual timeout:   %+v", defaultAddr, s.addr)
		}
	})

	t.Run("without tls", func(t *testing.T) {
		b, logger := setupLogger(t)
		NewServer(WithServerLogger(logger))

		if len(b.Buffer) != 1 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
		expected := "level=info msg=\"http server created\" addr=:8443\n"
		if b.Buffer[0] != expected {
			t.Errorf("\nExpected: %sActual:   %s", expected, b.Buffer[0])
		}
	})

	t.Run("with tls enabled", func(t *testing.T) {
		b, logger := setupLogger(t)
		cfg := mock.Server{
			Cert_: "../test/insecure/cert.pem",
			Key_:  "../test/insecure/key.pem",
		}
		tlsConfig, err := config.NewTLS(cfg)
		if err != nil {
			t.Errorf("Failed to create tls config")
		}

		NewServer(WithServerTLS(tlsConfig), WithServerLogger(logger))
		if len(b.Buffer) != 1 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
		expected := "level=info msg=\"http server created\" addr=:8443 tls=enabled\n"
		if b.Buffer[0] != expected {
			t.Errorf("\nExpected: %sActual:   %s", expected, b.Buffer[0])
		}
	})

	t.Run("with mtls enabled", func(t *testing.T) {
		b, logger := setupLogger(t)
		cfg := mock.Server{
			Cert_:   "../test/insecure/cert.pem",
			Key_:    "../test/insecure/key.pem",
			CACert_: "../test/insecure/rootCA.pem",
		}
		tlsConfig, err := config.NewTLS(cfg)
		if err != nil {
			t.Errorf("Failed to create mtls config")
		}

		NewServer(WithServerTLS(tlsConfig), WithServerLogger(logger))
		if len(b.Buffer) != 1 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
		expected := "level=info msg=\"http server created\" addr=:8443 mtls=enabled\n"
		if b.Buffer[0] != expected {
			t.Errorf("\nExpected: %sActual:   %s", expected, b.Buffer[0])
		}
	})
}

func TestServerRegister(t *testing.T) {
	b, logger := setupLogger(t)
	m := &mockService{}

	s := NewServer(WithServerLogger(logger))
	s.Register(m)
	if !m.registerCalled {
		t.Errorf("Failed to call register")
	}

	if len(b.Buffer) != 2 {
		t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
	}
	expected := "level=info msg=\"http server created\" addr=:8443\n"
	if b.Buffer[0] != expected {
		t.Errorf("\nExpected: %sActual:   %s", expected, b.Buffer[0])
	}
	expected = "level=info msg=\"service registered\" service=mockService\n"
	if b.Buffer[1] != expected {
		t.Errorf("\nExpected: %sActual:   %s", expected, b.Buffer[1])
	}

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	rw := httptest.NewRecorder()
	s.router.ServeHTTP(rw, req)

	if !m.handlerCalled {
		t.Errorf("Failed to call handler")
	}
	if rw.Code != http.StatusOK {
		t.Errorf("Unexpected status: %+v", rw.Code)
	}
}

func TestServerServe(t *testing.T) {
	t.Run("listen failure", func(t *testing.T) {
		s := NewServer(WithServerAddr(":42"))
		err := s.Serve()
		if err == nil {
			t.Errorf("Expected an error")
		}
		expected := "Failed to listen on :42"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("\nExpected error: %s\nActual error:   %s", err.Error(), expected)
		}
	})

	t.Run("serve no-tls", func(t *testing.T) {
		m := &mockService{}

		s := NewServer(WithServerAddr(":4242"))
		s.Register(m)

		result := make(chan error, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { wg.Done(); result <- s.Serve() }()
		wg.Wait()

		s.Shutdown()
		err := <-result

		if err != nil {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("serve mtls", func(t *testing.T) {
		cfg := mock.Server{
			Cert_:   "../test/insecure/cert.pem",
			Key_:    "../test/insecure/key.pem",
			CACert_: "../test/insecure/rootCA.pem",
		}
		tlsConfig, err := config.NewTLS(cfg)
		if err != nil {
			t.Errorf("Failed to create mtls config")
		}
		m := &mockService{}

		s := NewServer(WithServerTLS(tlsConfig), WithServerAddr(":2424"))
		s.Register(m)

		result := make(chan error, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { wg.Done(); result <- s.Serve() }()
		wg.Wait()

		s.Shutdown()
		err = <-result

		if err != nil {
			t.Errorf("Unexpected err: %s", err)
		}
	})
}

func TestServerRecoverPanics(t *testing.T) {
	b, logger := setupLogger(t)
	m := &mockService{}

	s := NewServer(WithServerLogger(logger))
	s.Register(m)

	req := httptest.NewRequest(http.MethodPut, "/panic", nil)
	rw := httptest.NewRecorder()
	s.router.ServeHTTP(rw, req)

	if !m.panicCalled {
		t.Errorf("Failed to call panic")
	}
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected status: %+v", rw.Code)
	}

	n := len(b.Buffer)
	expected := "level=error msg=\"recovering from a panic\""
	if !strings.Contains(b.Buffer[n-2], expected) {
		t.Errorf("\nExpected: %s\nActual:   %s", expected, b.Buffer[n-2])
	}
	expected = "level=error stacktrace="
	if !strings.Contains(b.Buffer[n-1], expected) {
		t.Errorf("\nExpected: %s\nActual:   %s", expected, b.Buffer[n-1])
	}
}

func TestServerRequestLogging(t *testing.T) {
	b, logger := setupLogger(t)
	m := &mockService{}

	s := NewServer(WithServerLogger(logger))
	s.Register(m)

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	rw := httptest.NewRecorder()
	s.router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Unexpected status: %+v", rw.Code)
	}

	n := len(b.Buffer)
	expected := `level=info method=GET url=/foo msg="request start"`
	if !strings.Contains(b.Buffer[n-3], expected) {
		t.Errorf("\nExpected: %s\nActual:   %s", expected, b.Buffer[n-3])
	}
	expected = `level=debug method=GET url=/foo msg="handler called"`
	if !strings.Contains(b.Buffer[n-2], expected) {
		t.Errorf("\nExpected: %s\nActual:   %s", expected, b.Buffer[n-2])
	}
	expected = `level=info method=GET url=/foo msg="request stop"`
	if !strings.Contains(b.Buffer[n-1], expected) {
		t.Errorf("\nExpected: %s\nActual:   %s", expected, b.Buffer[n-1])
	}
}

func setupLogger(t *testing.T) (*test.StringBuffer, log.Logger) {
	t.Helper()

	b := test.NewStringBuffer()
	logger, err := log.New(
		log.WithLevel(log.LevelDebug),
		log.WithFormat(log.FormatLogfmt),
		log.WithOutput(b),
		log.WithoutTimestamp(),
	)
	if err != nil {
		t.Fatal("failed to create logger")
	}

	return b, logger
}

type (
	mockService struct {
		registerCalled, handlerCalled, panicCalled bool
	}
)

var _ Service = (*mockService)(nil)

func (m *mockService) Register(r *mux.Router) {
	m.registerCalled = true
	r.HandleFunc("/foo", m.handler).Methods(http.MethodGet)
	r.HandleFunc("/panic", m.boom).Methods(http.MethodPut)
}

func (m *mockService) Fields() []interface{} {
	return []interface{}{
		"service", "mockService",
	}
}

func (m *mockService) handler(w http.ResponseWriter, r *http.Request) {
	m.handlerCalled = true
	logger := log.LoggerFromContext(r.Context())
	logger.Debug("msg", "handler called")
	w.WriteHeader(http.StatusOK)
}

func (m *mockService) boom(http.ResponseWriter, *http.Request) {
	m.panicCalled = true
	panic("boom")
}
