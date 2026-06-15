package telnet_test

import (
	"context"
	"testing"
	"time"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/log"
	"arcadium.dev/core/require"
	"arcadium.dev/core/telnet"
	"github.com/rs/zerolog"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	b := log.NewStringBuffer()
	logger, err := log.New(
		log.WithOutput(b),
		log.WithLevel(zerolog.DebugLevel),
		log.WithLevelFieldName("severity"),
		log.WithoutTimestamp(),
		log.AsDefault(),
	)
	require.Nil(t, err)

	tests := []struct {
		name   string
		opts   []telnet.ServerOption
		verify func(*testing.T, *telnet.Server)
	}{
		{
			name: "test defaults",
			verify: func(t *testing.T, s *telnet.Server) {
				require.NotNil(t, s)
				assert.Equal(t, s.Addr(), telnet.DefaultAddr)
			},
		},
		{
			name: "test WithServerAddress option",
			opts: []telnet.ServerOption{
				telnet.WithServerAddress(":2323"),
			},
			verify: func(t *testing.T, s *telnet.Server) {
				require.NotNil(t, s)
				assert.Equal(t, s.Addr(), ":2323")
			},
		},
		{
			name: "test WithServerLogger option",
			opts: []telnet.ServerOption{
				telnet.WithServerLogger(logger),
			},
			verify: func(t *testing.T, s *telnet.Server) {
				assert.NotNil(t, s)
				assert.Equal(t, s.Addr(), telnet.DefaultAddr)
				assert.Equal(t, s.Logger(), logger)

				require.Equal(t, b.Len(), 1)
				assert.Equal(t, b.Index(0), `{"severity":"info","address":":23","message":"telnet server created"}`+"\n")
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s := telnet.NewServer(test.opts...)
			test.verify(t, s)
		})
	}
}

func TestRegister(t *testing.T) {
	t.Parallel()

	s := telnet.NewServer()
	require.NotNil(t, s)

	m := &mockService{}
	s.Register(m)
	assert.True(t, m.registerCalled)

	mw := &mockMiddleware{}
	s.Middleware(mw.mock)

	s.Handle(&telnet.Session{})
	assert.True(t, m.handlerCalled)
	assert.True(t, mw.called)
	assert.True(t, mw.whenHandled.Before(m.whenHandled))
}

func TestServe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []telnet.ServerOption
		ctx     func() context.Context
		service *mockService
		verify  func(*testing.T, *mockService, error)
	}{
		{
			name: "listen failure",
			opts: []telnet.ServerOption{
				telnet.WithServerAddress(":-42"),
			},
			ctx:     func() context.Context { return context.Background() },
			service: &mockService{},
			verify: func(t *testing.T, service *mockService, err error) {
				assert.Error(t, err, "listen tcp: address -42: invalid port")
				assert.True(t, service.shutdownCalled)
			},
		},
		{
			name: "cancelled context",
			opts: []telnet.ServerOption{
				telnet.WithServerAddress(":2323"),
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			service: &mockService{},
			verify: func(t *testing.T, service *mockService, err error) {
				assert.Nil(t, err)
				assert.True(t, service.shutdownCalled)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := telnet.NewServer(test.opts...)
			s.Register(test.service)

			mw := &mockMiddleware{}
			s.Middleware(mw.mock)

			result := make(chan error, 1)
			go func() {
				result <- s.Serve()
			}()

			var err error
			select {
			case <-test.ctx().Done():
			case err = <-result:
			}
			s.Shutdown(context.Background())
			test.verify(t, test.service, err)
		})
	}
}

type (
	mockService struct {
		registerCalled, handlerCalled, shutdownCalled bool
		whenHandled                                   time.Time
	}
)

var _ telnet.Service = (*mockService)(nil)

func (m *mockService) Register(s *telnet.Server) {
	m.registerCalled = true
}

func (m *mockService) Name() string {
	return "mockService"
}

func (m *mockService) ServeTELNET(s *telnet.Session) {
	m.handlerCalled = true
	m.whenHandled = time.Now()
}

func (m *mockService) Shutdown(_ context.Context) {
	m.shutdownCalled = true
}

type (
	mockMiddleware struct {
		called      bool
		whenHandled time.Time
	}
)

func (m *mockMiddleware) mock(next telnet.Handler) telnet.Handler {
	return telnet.HandlerFunc(func(session *telnet.Session) {
		m.called = true
		m.whenHandled = time.Now()
		next.ServeTELNET(session)
	})
}
