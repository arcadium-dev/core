package rest_test

import (
	"context"
	"errors"
	"testing"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/http/server"
	"arcadium.dev/core/log"
	"arcadium.dev/core/require"
	"arcadium.dev/core/rest"
	"github.com/rs/zerolog"
)

func TestServerInit(t *testing.T) {
	t.Run("new config failure", func(t *testing.T) {
		s, b := setup(t)

		s.C.NewConfig = func() (rest.Config, error) {
			return rest.Config{}, errors.New("new config failure")
		}

		err := s.Init()

		assert.Error(t, err, "failed to load config: new config failure")
		require.Equal(t, b.Len(), 1)
		assert.Contains(t, b.Index(0), "failed to load config: new config failure")
	})

	t.Run("new logger failure", func(t *testing.T) {
		s, b := setup(t)

		s.C.NewLogger = func(rest.Config) (*zerolog.Logger, error) {
			return nil, errors.New("new logger failure")
		}

		err := s.Init()

		assert.Error(t, err, "failed to create logger: new logger failure")
		require.Equal(t, b.Len(), 1)
		assert.Contains(t, b.Index(0), "failed to create logger: new logger failure")
	})
}

func TestServerStart(t *testing.T) {
	t.Run("new http server failure", func(t *testing.T) {
		s, b := setup(t)

		s.C.NewHTTPServer = func(context.Context, rest.Config) (*server.Server, error) {
			return nil, errors.New("new http server failure")
		}

		assert.Nil(t, s.Init())
		err := s.Start()

		assert.Error(t, err, "failed to create http server: new http server failure")
		require.Equal(t, b.Len(), 2)
		assert.Contains(t, b.Index(1), `"error":"new http server failure"`)
	})

	t.Run("success", func(t *testing.T) {
		s, _ := setup(t)

		assert.Nil(t, s.Init())

		r := make(chan error, 1)
		go func() { r <- s.Start() }()
		s.Stop()
		err := <-r

		assert.Nil(t, err)
	})
}

func setup(t *testing.T) (*rest.Server, *log.StringBuffer) {
	s := rest.NewServer("version", "branch", "commit", "date")

	b := log.NewStringBuffer()
	s.Stdout = b
	s.Stderr = b

	t.Setenv("DSN", "root:password@tcp(mariadb:3306)/dbname")
	t.Setenv("SERVER_ADDR", ":8443")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("TLS_CERT", "")
	t.Setenv("TLS_KEY", "")
	t.Setenv("MTLS_ENABLED", "false")
	t.Setenv("ALLOWED_ORIGINS", "https://*.arcadium.dev")
	t.Setenv("ALLOWED_METHODS", "GET")
	t.Setenv("ALLOWED_HEADERS", "*")

	return s, b
}
