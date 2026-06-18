package telnet_test

import (
	"context"
	"testing"

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

	server.Handle(&telnet.Session{})

	assert.True(t, service.called)
	n := b.Len()
	assert.Contains(t, b.Index(n-2), `{"severity":"debug","message":"session start"`)
	assert.Contains(t, b.Index(n-1), `{"severity":"debug","message":"session complete"`)
}

type SessionService struct {
	called bool
}

func (s SessionService) Name() string                         { return "session" }
func (s *SessionService) ServeTELNET(session *telnet.Session) { s.called = true }
func (s SessionService) Shutdown(context.Context)             {}
