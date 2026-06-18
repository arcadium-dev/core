package telnet_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/log"
	"arcadium.dev/core/require"
	"arcadium.dev/core/telnet"
)

func TestRecoveryMiddleware_Recover(t *testing.T) {
	b := log.NewStringBuffer()
	logger, err := log.New(
		log.WithOutput(b),
		log.WithLevel(zerolog.DebugLevel),
		log.WithLevelFieldName("severity"),
		log.WithoutTimestamp(),
		log.AsDefault(),
	)
	require.Nil(t, err)

	server := telnet.NewServer()
	server.Register(BoomService{})
	server.Middleware(telnet.RecoveryMiddleware{Logger: logger}.Recover)

	server.Handle(&telnet.Session{})

	n := b.Len()
	assert.Contains(t, b.Index(n-2), `{"severity":"error","message":"recovering from a panic"}`)
	assert.Contains(t, b.Index(n-1), `{"severity":"error","message":"stacktrace:`)
}

type BoomService struct{}

func (b BoomService) Name() string                        { return "boom" }
func (b BoomService) ServeTELNET(session *telnet.Session) { panic("boom") }
func (b BoomService) Shutdown(context.Context)            {}

func TestSessionMiddleware_Session(t *testing.T) {
	b := log.NewStringBuffer()
	logger, err := log.New(
		log.WithOutput(b),
		log.WithLevel(zerolog.DebugLevel),
		log.WithLevelFieldName("severity"),
		log.WithoutTimestamp(),
		log.AsDefault(),
	)
	require.Nil(t, err)

	server := telnet.NewServer()
	service := SessionService{}

	server.Register(&service)
	server.Middleware(
		telnet.RecoveryMiddleware{Logger: logger}.Recover, // to test the path when a session doesn't panic
		telnet.SessionMiddleware{Logger: logger}.Session,
	)

	server.Handle(&telnet.Session{
		Conn: mockConn{addr: mockAddr("myaddress:1234")},
	})

	assert.True(t, service.called)
	n := b.Len()
	assert.Equal(t, b.Index(n-2), `{"severity":"debug","remote addr":"myaddress:1234","message":"session start"}`+"\n")
	assert.Equal(t, b.Index(n-1), `{"severity":"debug","remote addr":"myaddress:1234","message":"session complete"}`+"\n")
}

type SessionService struct {
	called bool
}

func (s SessionService) Name() string                         { return "session" }
func (s *SessionService) ServeTELNET(session *telnet.Session) { s.called = true }
func (s SessionService) Shutdown(context.Context)             {}

type mockConn struct {
	addr mockAddr
}

func (m mockConn) RemoteAddr() net.Addr { return m.addr }

func (m mockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m mockConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (m mockConn) Close() error                       { return nil }
func (m mockConn) LocalAddr() net.Addr                { return nil }
func (m mockConn) SetDeadline(t time.Time) error      { return nil }
func (m mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m mockConn) SetWriteDeadline(t time.Time) error { return nil }

type mockAddr string

func (m mockAddr) Network() string { return "" }
func (m mockAddr) String() string  { return string(m) }
