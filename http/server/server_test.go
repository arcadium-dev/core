package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/log"
	"arcadium.dev/core/require"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func TestServerNew(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		s := NewServer(context.Background())

		assert.Equal(t, s.addr, defaultAddr)
		assert.Equal(t, s.shutdownTimeout, defaultShutdownTimeout)
	})

	t.Run("without tls", func(t *testing.T) {
		ctx, b := log.SetupTestLogging(t)
		NewServer(ctx)

		require.Equal(t, b.Len(), 2)
		assert.Equal(t, b.Index(1), `{"severity":"info","message":"http server created, address ':8443'"}`+"\n")
	})

	t.Run("with tls enabled", func(t *testing.T) {
		ctx, b := log.SetupTestLogging(t)
		cfg := setupTLS(t, "./test/insecure_cert.pem", "./test/insecure_key.pem")

		NewServer(ctx, WithTLS(cfg))

		require.Equal(t, b.Len(), 2)
		assert.Equal(t, b.Index(1), `{"severity":"info","message":"http server created, address ':8443', tls: enabled"}`+"\n")
	})
}

func TestServerRegister(t *testing.T) {
	ctx, b := log.SetupTestLogging(t)
	m := &mockService{}
	s := NewServer(ctx)

	s.Register(m)

	if !m.registerCalled {
		t.Errorf("Failed to call register")
	}

	require.Equal(t, b.Len(), 3)
	assert.Equal(t, b.Index(0), `{"severity":"info","message":"cors allow all"}`+"\n")
	assert.Equal(t, b.Index(1), `{"severity":"info","message":"http server created, address ':8443'"}`+"\n")
	assert.Equal(t, b.Index(2), `{"severity":"info","message":"http service registered: mockService"}`+"\n")

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	rw := httptest.NewRecorder()
	s.router.ServeHTTP(rw, req)

	assert.True(t, m.handlerCalled)
	assert.Equal(t, rw.Code, http.StatusOK)
}

func TestServerCORS(t *testing.T) {
	t.Run("preflight abort - origin", func(t *testing.T) {
		s := NewServer(context.Background(), WithCORS(&cors.Options{
			AllowedOrigins: []string{"https://*.arcadium.dev"},
			AllowedMethods: []string{"GET"},
			AllowedHeaders: []string{"X-Requested-With", "Content-Type"},
		}))

		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Origin", "http://www.arcadium.dev")
		r.Header.Set("Access-Control-Request-Method", "GET")
		r.Header.Set("Access-Control-Request-Headers", "X-Requested-With,Content-Type")

		s.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusNoContent)

		acao := w.Header().Get("Access-Control-Allow-Origin")
		assert.Equal(t, acao, "")

		acam := w.Header().Get("Access-Control-Allow-Methods")
		assert.Equal(t, acam, "")

		acah := w.Header().Get("Access-Control-Allow-Headers")
		assert.Equal(t, acah, "")
	})

	t.Run("preflight abort - method", func(t *testing.T) {
		s := NewServer(context.Background(), WithCORS(&cors.Options{
			AllowedOrigins: []string{"https://*.arcadium.dev"},
			AllowedMethods: []string{"GET"},
			AllowedHeaders: []string{"X-Requested-With", "Content-Type"},
		}))

		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Origin", "https://arcade.arcadium.dev")
		r.Header.Set("Access-Control-Request-Method", "PUT")
		r.Header.Set("Access-Control-Request-Headers", "X-Requested-With,Content-Type")

		s.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusNoContent)

		acao := w.Header().Get("Access-Control-Allow-Origin")
		assert.Equal(t, acao, "")

		acam := w.Header().Get("Access-Control-Allow-Methods")
		assert.Equal(t, acam, "")

		acah := w.Header().Get("Access-Control-Allow-Headers")
		assert.Equal(t, acah, "")
	})

	t.Run("preflight abort - header", func(t *testing.T) {
		s := NewServer(context.Background(), WithCORS(&cors.Options{
			AllowedOrigins: []string{"https://*.arcadium.dev"},
			AllowedMethods: []string{"GET"},
			AllowedHeaders: []string{"X-Requested-With", "Content-Type"},
		}))

		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Origin", "https://arcade.arcadium.dev")
		r.Header.Set("Access-Control-Request-Method", "GET")
		r.Header.Set("Access-Control-Request-Headers", "X-Requested-With,Content-Type,x-okta-user-agent-extended")

		s.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusNoContent)

		acao := w.Header().Get("Access-Control-Allow-Origin")
		assert.Equal(t, acao, "")

		acam := w.Header().Get("Access-Control-Allow-Methods")
		assert.Equal(t, acam, "")

		acah := w.Header().Get("Access-Control-Allow-Headers")
		assert.Equal(t, acah, "")
	})

	t.Run("success - default cors", func(t *testing.T) {
		ctx, b := log.SetupTestLogging(t)
		s := NewServer(ctx)

		require.Equal(t, b.Len(), 2)
		assert.Equal(t, b.Index(0), `{"severity":"info","message":"cors allow all"}`+"\n")
		assert.Equal(t, b.Index(1), `{"severity":"info","message":"http server created, address ':8443'"}`+"\n")

		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Origin", "http://www.arcadium.dev")
		r.Header.Set("Access-Control-Request-Method", "GET")
		r.Header.Set("Access-Control-Request-Headers", "X-Requested-With,Content-Type,x-okta-user-agent-extended")

		s.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusNoContent)

		acao := w.Header().Get("Access-Control-Allow-Origin")
		assert.Equal(t, acao, "*")

		acam := w.Header().Get("Access-Control-Allow-Methods")
		assert.Equal(t, acam, http.MethodGet)

		acah := w.Header().Get("Access-Control-Allow-Headers")
		assert.Equal(t, acah, "X-Requested-With, Content-Type, X-Okta-User-Agent-Extended")
	})

	t.Run("success - custom cors", func(t *testing.T) {
		ctx, b := log.SetupTestLogging(t)
		s := NewServer(ctx, WithCORS(&cors.Options{
			AllowedOrigins: []string{"https://*.arcadium.dev"},
			AllowedMethods: []string{"GET"},
			AllowedHeaders: []string{"X-Requested-With", "Content-Type"},
		}))

		require.Equal(t, b.Len(), 4)
		assert.Equal(t, b.Index(0), `{"severity":"info","message":"cors allowed origins: [\"https://*.arcadium.dev\"]"}`+"\n")
		assert.Equal(t, b.Index(1), `{"severity":"info","message":"cors allowed methods: [\"GET\"]"}`+"\n")
		assert.Equal(t, b.Index(2), `{"severity":"info","message":"cors allowed headers: [\"X-Requested-With\" \"Content-Type\"]"}`+"\n")
		assert.Equal(t, b.Index(3), `{"severity":"info","message":"http server created, address ':8443'"}`+"\n")

		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Origin", "https://arcade.arcadium.dev")
		r.Header.Set("Access-Control-Request-Method", "GET")
		r.Header.Set("Access-Control-Request-Headers", "X-Requested-With,Content-Type")

		s.server.Handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusNoContent)

		acao := w.Header().Get("Access-Control-Allow-Origin")
		assert.Equal(t, acao, "https://arcade.arcadium.dev")

		acam := w.Header().Get("Access-Control-Allow-Methods")
		assert.Equal(t, acam, http.MethodGet)

		acah := w.Header().Get("Access-Control-Allow-Headers")
		assert.Equal(t, acah, "X-Requested-With, Content-Type")
	})
}

func TestServerServe(t *testing.T) {
	t.Run("listen failure", func(t *testing.T) {
		s := NewServer(context.Background(), WithAddr(":-42"))
		err := s.Serve()
		assert.Contains(t, err.Error(), "failed to listen on address ':-42'")
	})

	t.Run("serve", func(t *testing.T) {
		m := &mockService{}
		s := NewServer(context.Background(), WithAddr(":4242"))

		s.Register(m)

		assert.Equal(t, len(s.services), 1)
		assert.Equal(t, s.services[0].Name(), "mockService")

		result := make(chan error, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { wg.Done(); result <- s.Serve() }()
		wg.Wait()

		s.Shutdown()
		err := <-result

		assert.Nil(t, err)
	})

	t.Run("serve with middleware", func(t *testing.T) {
		m := &mockService{}

		s := NewServer(context.Background(), WithAddr(":4242"))
		s.Middleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					r := recover()
					actual, ok := r.(string)
					require.True(t, ok)
					assert.Equal(t, actual, "boom")
				}()
				next.ServeHTTP(w, r)
			})
		})
		s.Register(m)

		assert.Equal(t, len(s.services), 1)
		assert.Equal(t, s.services[0].Name(), "mockService")

		result := make(chan error, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { wg.Done(); result <- s.Serve() }()
		wg.Wait()

		req := httptest.NewRequest(http.MethodGet, "/boom", nil)
		rw := httptest.NewRecorder()
		s.router.ServeHTTP(rw, req)

		s.Shutdown()
		err := <-result

		assert.True(t, m.shutdownCalled())
		assert.Nil(t, err)
	})

	t.Run("serve tls", func(t *testing.T) {
		tlsConfig := setupTLS(t, "./test/insecure_cert.pem", "./test/insecure_key.pem")
		m := &mockService{}
		s := NewServer(context.Background(), WithTLS(tlsConfig), WithAddr(":2424"))
		s.Register(m)

		require.Equal(t, len(s.services), 1)
		assert.Equal(t, s.services[0].Name(), "mockService")

		result := make(chan error, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { wg.Done(); result <- s.Serve() }()
		wg.Wait()

		s.Shutdown()
		err := <-result

		assert.True(t, m.shutdownCalled())
		assert.Nil(t, err)
	})
}

func setupTLS(t *testing.T, cert, key string) *tls.Config {
	t.Helper()

	certificate, err := tls.LoadX509KeyPair(cert, key)
	require.Nil(t, err)

	cfg := &tls.Config{}
	cfg.Certificates = append(cfg.Certificates, certificate)

	return cfg
}

type (
	mockService struct {
		registerCalled, handlerCalled, panicCalled bool

		mu       sync.RWMutex
		shutdown bool
	}
)

var _ Service = (*mockService)(nil)

func (m *mockService) Register(r *mux.Router) {
	m.registerCalled = true
	r.HandleFunc("/foo", m.handler).Methods(http.MethodGet)
	r.HandleFunc("/panic", m.boom).Methods(http.MethodPut)
}

func (m *mockService) Name() string {
	return "mockService"
}

func (m *mockService) Shutdown() {
	m.mu.Lock()
	m.shutdown = true
	m.mu.Unlock()
}

func (m *mockService) handler(w http.ResponseWriter, r *http.Request) {
	m.handlerCalled = true
	w.WriteHeader(http.StatusOK)
}

func (m *mockService) boom(http.ResponseWriter, *http.Request) {
	m.panicCalled = true
	panic("boom")
}

func (m *mockService) shutdownCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.shutdown
}
