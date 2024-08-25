package ssh_test

import (
	"context"
	"errors"
	"testing"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/http/server"
	"arcadium.dev/core/log"
	"arcadium.dev/core/require"
	"arcadium.dev/core/ssh"
	"github.com/rs/zerolog"
)

func TestServerInit(t *testing.T) {
	t.Run("new config failure", func(t *testing.T) {
		s, b := setup(t)

		s.C.NewConfig = func(...string) (ssh.Config, error) {
			return ssh.Config{}, errors.New("new config failure")
		}

		err := s.Init()

		assert.Error(t, err, "failed to load config: new config failure")
		require.Equal(t, b.Len(), 1)
		assert.Contains(t, b.Index(0), "failed to load config: new config failure")
	})

	t.Run("new logger failure", func(t *testing.T) {
		s, b := setup(t)

		s.C.NewLogger = func(ssh.Config) (*zerolog.Logger, error) {
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

		s.C.NewHTTPServer = func(context.Context, ssh.Config) (*server.Server, error) {
			return nil, errors.New("new http server failure")
		}

		assert.Nil(t, s.Init())
		err := s.Start()

		assert.Error(t, err, "failed to create http server: new http server failure")
		require.Equal(t, b.Len(), 2)
		assert.Contains(t, b.Index(1), `"error":"new http server failure"`)
	})

	t.Run("success", func(t *testing.T) {
		s, _ := setup(t, "TEST_")

		assert.Nil(t, s.Init("test"))

		r := make(chan error, 1)
		go func() { r <- s.Start() }()
		s.Stop()
		err := <-r

		assert.Nil(t, err)
	})
}

func setup(t *testing.T, prefix ...string) (*ssh.Server, *log.StringBuffer) {
	s := ssh.NewServer("version", "branch", "commit", "date")

	b := log.NewStringBuffer()
	s.Stdout = b

	p := ""
	if len(prefix) == 1 {
		p = prefix[0]
	}

	t.Setenv(p+"HTTP_SERVER_ADDR", ":8443")
	t.Setenv(p+"LOG_LEVEL", "debug")
	t.Setenv(p+"TLS_CERT", "")
	t.Setenv(p+"TLS_KEY", "")
	t.Setenv(p+"MTLS_ENABLED", "false")
	t.Setenv(p+"ALLOWED_ORIGINS", "https://*.arcadium.dev")
	t.Setenv(p+"ALLOWED_METHODS", "GET")
	t.Setenv(p+"ALLOWED_HEADERS", "*")

	return s, b
}
