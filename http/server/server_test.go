package server_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/http/server"
	"arcadium.dev/core/log"
	"arcadium.dev/core/require"
)

func TestServer_New(t *testing.T) {
	tests := []struct {
		name   string
		opts   []server.Option
		verify func(*testing.T, *server.Server, error)
	}{
		// Test WithAddr option.
		{
			name: "with addr",
			opts: []server.Option{server.WithAddr("10.11.12.13:4201")},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				assert.Equal(t, s.Addr(), "10.11.12.13:4201")
			},
		},

		// Test WithTLS options.
		{
			name: "with tls config",
			opts: []server.Option{server.WithTLSConfig(&tls.Config{})},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s)
				assert.NotNil(t, s.TLSConfig())
			},
		},
		{
			name: "with tls config, and tls properties",
			opts: []server.Option{
				server.WithTLSConfig(&tls.Config{}),
				server.WithTLSCert("./test/bad_cert.pem", "./test/bad_key.pem"),
			},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s)
				assert.NotNil(t, s.TLSConfig())
			},
		},
		{
			name: "with tls cert, without tls key failure",
			opts: []server.Option{server.WithTLSCert("./test/insecure_cert.pem", "")},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, s)
				assert.Error(t, err, `the tls key must be defined with then tls cert is defined`)
			},
		},
		{
			name: "without tls cert, with tls key failure",
			opts: []server.Option{server.WithTLSCert("", "./test/insecure_key.pem")},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, s)
				assert.Error(t, err, `the tls cert must be defined with then tls key is defined`)
			},
		},
		{
			name: "with bad tls cert, tls key failure",
			opts: []server.Option{server.WithTLSCert("./test/bad_cert.pem", "./test/bad_key.pem")},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, s)
				assert.Error(t, err, `failed to load TLS certificate: open ./test/bad_cert.pem: no such file or directory`)
			},
		},
		{
			name: "with tls cert, with tls key, without client ca cert",
			opts: []server.Option{server.WithTLSCert("./test/insecure_cert.pem", "./test/insecure_key.pem")},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s.TLSConfig())
				assert.Equal(t, len(s.TLSConfig().Certificates), 1)
				assert.Nil(t, s.TLSConfig().ClientCAs)
				assert.Equal(t, s.TLSConfig().ClientAuth, tls.NoClientCert)
			},
		},
		{
			name: "with tls cert, with tls key, with client ca cert failure",
			opts: []server.Option{
				server.WithTLSCert("./test/insecure_cert.pem", "./test/insecure_key.pem"),
				server.WithTLSClientCACert("./test/bad_rootCA.pem"),
			},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, s)
				assert.Error(t, err, `failed to load the TLS client CA certificate: open ./test/bad_rootCA.pem: no such file or directory`)
			},
		},
		{
			name: "with tls cert, with tls key failure, with client ca cert",
			opts: []server.Option{
				server.WithTLSCert("./test/insecure_cert.pem", "./test/insecure_key.pem"),
				server.WithTLSClientCACert("./test/insecure_rootCA.pem"),
			},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s.TLSConfig())
				assert.Equal(t, len(s.TLSConfig().Certificates), 1)
				assert.NotNil(t, s.TLSConfig().ClientCAs)
				assert.Equal(t, s.TLSConfig().ClientAuth, tls.NoClientCert)
			},
		},
		{
			name: "with tls cert, with tls key failure, with client ca cert, with mtls enabled",
			opts: []server.Option{
				server.WithTLSCert("./test/insecure_cert.pem", "./test/insecure_key.pem"),
				server.WithTLSClientCACert("./test/insecure_rootCA.pem"),
				server.WithMTLSEnabled(true),
			},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s.TLSConfig())
				assert.Equal(t, len(s.TLSConfig().Certificates), 1)
				assert.NotNil(t, s.TLSConfig().ClientCAs)
				assert.Equal(t, s.TLSConfig().ClientAuth, tls.RequireAndVerifyClientCert)
			},
		},

		// Test WithCORSOptions option.
		{
			name: "with cors options",
			opts: []server.Option{server.WithCORSOptions(&cors.Options{})},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s)
				assert.NotNil(t, s.CORSOptions())
			},
		},
		{
			name: "with cors options and cors properties failure",
			opts: []server.Option{
				server.WithCORSOptions(&cors.Options{}),
				server.WithCORSAllowedOrigins([]string{"*.arcadium.dev"}),
			},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s)
				assert.NotNil(t, s.CORSOptions())
			},
		},
		{
			name: "with cors properties",
			opts: []server.Option{
				server.WithCORSAllowedOrigins([]string{"*"}),
				server.WithCORSAllowedMethods([]string{"GET"}),
				server.WithCORSAllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
			},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s)
				opts := s.CORSOptions()
				require.NotNil(t, opts)
				assert.Compare(t, opts.AllowedOrigins, []string{"*"})
				assert.Compare(t, opts.AllowedMethods, []string{"GET"})
				assert.Compare(t, opts.AllowedHeaders, []string{"X-Requested-With", "Content-Type"})
				assert.False(t, opts.AllowCredentials)
			},
		},
		{
			name: "with cors properties - allow credentials",
			opts: []server.Option{
				server.WithCORSAllowedOrigins([]string{"arcadium.dev"}),
				server.WithCORSAllowedMethods([]string{"GET"}),
				server.WithCORSAllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
			},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s)
				opts := s.CORSOptions()
				require.NotNil(t, opts)
				assert.Compare(t, opts.AllowedOrigins, []string{"arcadium.dev"})
				assert.Compare(t, opts.AllowedMethods, []string{"GET"})
				assert.Compare(t, opts.AllowedHeaders, []string{"X-Requested-With", "Content-Type"})
				assert.True(t, opts.AllowCredentials)
			},
		},

		// Test timeouts.
		{
			name: "defaults",
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				assert.Equal(t, s.ReadTimeout(), 5*time.Second)
				assert.Equal(t, s.WriteTimeout(), 10*time.Second)
				assert.Equal(t, s.ShutdownTimeout(), 10*time.Second)
			},
		},
		{
			name: "with timeouts",
			opts: []server.Option{
				server.WithReadTimeout(7 * time.Second),
				server.WithWriteTimeout(13 * time.Second),
				server.WithShutdownTimeout(111 * time.Second),
			},
			verify: func(t *testing.T, s *server.Server, err error) {
				assert.Nil(t, err)
				assert.Equal(t, s.ReadTimeout(), 7*time.Second)
				assert.Equal(t, s.WriteTimeout(), 13*time.Second)
				assert.Equal(t, s.ShutdownTimeout(), 111*time.Second)
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			s, err := server.New(context.Background(), test.opts...)
			test.verify(t, s, err)
		})
	}
}

func TestServer_Register(t *testing.T) {
	ctx, b := log.SetupTestLogging(t)
	m := &mockService{}
	s, err := server.New(ctx)
	assert.Nil(t, err)

	s.Register(ctx, m)

	if !m.registerCalled {
		t.Errorf("Failed to call register")
	}

	require.Equal(t, b.Len(), 3)
	assert.Equal(t, b.Index(0), `{"severity":"info","message":"cors allow all"}`+"\n")
	assert.Equal(t, b.Index(1), `{"severity":"info","message":"http server created, address ':80'"}`+"\n")
	assert.Equal(t, b.Index(2), `{"severity":"info","message":"http service registered: mockService"}`+"\n")

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	rw := httptest.NewRecorder()
	s.Router().ServeHTTP(rw, req)

	assert.True(t, m.handlerCalled)
	assert.Equal(t, rw.Code, http.StatusOK)
}

func TestServer_CORS(t *testing.T) {
	tests := []struct {
		name    string
		opt     *cors.Options
		headers map[string]string
		verify  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "preflight abort - origin",
			opt: &cors.Options{
				AllowedOrigins: []string{"https://*.arcadium.dev"},
				AllowedMethods: []string{"GET"},
				AllowedHeaders: []string{"X-Requested-With", "Content-Type"},
			},
			headers: map[string]string{
				"Origin":                         "http://www.arcadium.dev",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Requested-With,Content-Type",
			},
			verify: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, w.Code, http.StatusNoContent)
				acao := w.Header().Get("Access-Control-Allow-Origin")
				assert.Equal(t, acao, "")
				acam := w.Header().Get("Access-Control-Allow-Methods")
				assert.Equal(t, acam, "")
				acah := w.Header().Get("Access-Control-Allow-Headers")
				assert.Equal(t, acah, "")
			},
		},
		{
			name: "preflight abort - method",
			opt: &cors.Options{
				AllowedOrigins: []string{"https://*.arcadium.dev"},
				AllowedMethods: []string{"GET"},
				AllowedHeaders: []string{"X-Requested-With", "Content-Type"},
			},
			headers: map[string]string{
				"Origin":                         "http://www.arcadium.dev",
				"Access-Control-Request-Method":  "PUT",
				"Access-Control-Request-Headers": "X-Requested-With,Content-Type",
			},
			verify: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, w.Code, http.StatusNoContent)
				acao := w.Header().Get("Access-Control-Allow-Origin")
				assert.Equal(t, acao, "")
				acam := w.Header().Get("Access-Control-Allow-Methods")
				assert.Equal(t, acam, "")
				acah := w.Header().Get("Access-Control-Allow-Headers")
				assert.Equal(t, acah, "")
			},
		},
		{
			name: "preflight abort - header",
			opt: &cors.Options{
				AllowedOrigins: []string{"https://*.arcadium.dev"},
				AllowedMethods: []string{"GET"},
				AllowedHeaders: []string{"X-Requested-With", "Content-Type"},
			},
			headers: map[string]string{
				"Origin":                         "https://arcade.arcadium.dev",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Requested-With,Content-Type,x-okta-user-agent-extended",
			},
			verify: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, w.Code, http.StatusNoContent)
				acao := w.Header().Get("Access-Control-Allow-Origin")
				assert.Equal(t, acao, "")
				acam := w.Header().Get("Access-Control-Allow-Methods")
				assert.Equal(t, acam, "")
				acah := w.Header().Get("Access-Control-Allow-Headers")
				assert.Equal(t, acah, "")
			},
		},
		{
			name: "success - default cors",
			headers: map[string]string{
				"Origin":                         "http://www.arcadium.dev",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Requested-With,Content-Type,X-Okta-User-Agent-Extended",
			},
			verify: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, w.Code, http.StatusNoContent)
				acao := w.Header().Get("Access-Control-Allow-Origin")
				assert.Equal(t, acao, "*")
				acam := w.Header().Get("Access-Control-Allow-Methods")
				assert.Equal(t, acam, http.MethodGet)
				acah := w.Header().Get("Access-Control-Allow-Headers")
				assert.Equal(t, acah, "X-Requested-With,Content-Type,X-Okta-User-Agent-Extended")
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			s, err := server.New(context.Background(), server.WithCORSOptions(test.opt))
			assert.Nil(t, err)
			require.NotNil(t, s)

			r := httptest.NewRequest(http.MethodOptions, "/", nil)
			w := httptest.NewRecorder()
			for key, value := range test.headers {
				r.Header.Set(key, value)
			}
			s.Server().Handler.ServeHTTP(w, r)

			require.NotNil(t, w)

			t.Logf("%+v", test)
			assert.NotNil(t, test.verify)
			test.verify(t, w)
		})
	}
}

func TestServer_Serve(t *testing.T) {
	t.Run("listen failure", func(t *testing.T) {
		s, err := server.New(context.Background(), server.WithAddr(":-42"))
		assert.Nil(t, err)
		err = s.Serve(context.Background())
		assert.Contains(t, err.Error(), "failed to listen on address ':-42'")
	})

	t.Run("serve", func(t *testing.T) {
		ctx := context.Background()
		m := &mockService{}
		s, err := server.New(context.Background(), server.WithAddr(":4242"))
		assert.Nil(t, err)

		s.Register(context.Background(), m)

		assert.Equal(t, s.Services().Len(), 1)
		assert.Equal(t, s.Services().Get(0).Name(), "mockService")

		result := make(chan error, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { wg.Done(); result <- s.Serve(ctx) }()
		wg.Wait()

		s.Shutdown(ctx)
		err = <-result

		assert.Nil(t, err)
	})

	t.Run("serve with middleware", func(t *testing.T) {
		ctx := context.Background()
		m := &mockService{}

		s, err := server.New(context.Background(), server.WithAddr(":4242"))
		assert.Nil(t, err)

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
		s.Register(ctx, m)

		assert.Equal(t, s.Services().Len(), 1)
		assert.Equal(t, s.Services().Get(0).Name(), "mockService")

		result := make(chan error, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { wg.Done(); result <- s.Serve(ctx) }()
		wg.Wait()

		req := httptest.NewRequest(http.MethodGet, "/boom", nil)
		rw := httptest.NewRecorder()
		s.Router().ServeHTTP(rw, req)

		s.Shutdown(ctx)
		err = <-result

		assert.True(t, m.shutdownCalled())
		assert.Nil(t, err)
	})

	t.Run("serve tls", func(t *testing.T) {
		ctx := context.Background()
		m := &mockService{}
		s, err := server.New(
			context.Background(),
			server.WithAddr(":2424"),
			server.WithTLSCert("./test/insecure_cert.pem", "./test/insecure_key.pem"),
		)
		assert.Nil(t, err)
		s.Register(ctx, m)

		require.Equal(t, s.Services().Len(), 1)
		assert.Equal(t, s.Services().Get(0).Name(), "mockService")

		result := make(chan error, 1)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { wg.Done(); result <- s.Serve(ctx) }()
		wg.Wait()

		s.Shutdown(ctx)
		err = <-result

		assert.True(t, m.shutdownCalled())
		assert.Nil(t, err)
	})
}

func TestServer_Name(t *testing.T) {
	s, err := server.New(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, s.Name(), "http server")
}

type (
	mockService struct {
		registerCalled, handlerCalled, panicCalled bool

		mu       sync.RWMutex
		shutdown bool
	}
)

var _ server.Service = (*mockService)(nil)

func (m *mockService) Register(r *mux.Router) {
	m.registerCalled = true
	r.HandleFunc("/foo", m.handler).Methods(http.MethodGet)
	r.HandleFunc("/panic", m.boom).Methods(http.MethodPut)
}

func (m *mockService) Name() string {
	return "mockService"
}

func (m *mockService) Shutdown(context.Context) {
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
