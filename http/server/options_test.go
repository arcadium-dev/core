package server

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/rs/cors"

	"arcadium.dev/core/assert"
)

func TestWithAddr(t *testing.T) {
	addr := "example.com:4201"
	s := &Server{}
	WithAddr(addr).apply(s)

	assert.Equal(t, s.addr, addr)
}

func TestWithTLS(t *testing.T) {
	s := &Server{
		server: &http.Server{},
	}
	cfg := &tls.Config{}
	WithTLS(cfg).apply(s)

	assert.NotNil(t, s.server.TLSConfig)
}

func TestWithCORS(t *testing.T) {
	corsOpts := &cors.Options{}
	s := &Server{}
	WithCORS(corsOpts).apply(s)

	assert.NotNil(t, s.corsOptions)
}

func TestWithReadTimeout(t *testing.T) {
	s := &Server{server: &http.Server{}}

	timeout := 50 * time.Second
	WithReadTimeout(timeout).apply(s)

	assert.Equal(t, s.server.ReadTimeout, timeout)
}

func TestWithWriteTimeout(t *testing.T) {
	s := &Server{server: &http.Server{}}

	timeout := 50 * time.Second
	WithWriteTimeout(timeout).apply(s)

	assert.Equal(t, s.server.WriteTimeout, timeout)
}

func TestWithShutdownTimeout(t *testing.T) {
	s := &Server{}
	timeout := 52 * time.Second
	WithShutdownTimeout(timeout).apply(s)

	assert.Equal(t, s.shutdownTimeout, timeout)
}
