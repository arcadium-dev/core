package server_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/http/server"
)

func TestNewConfig(t *testing.T) {
	for _, env := range []string{"DSN", "LOG_LEVEL"} {
		os.Unsetenv(env)
	}

	t.Run("test defaults", func(t *testing.T) {
		t.Setenv("TEST_SERVER_ADDR", ":8443")

		cfg, err := server.NewConfig("test")

		assert.Nil(t, err)
		assert.Equal(t, cfg.ServerAddr(), `:8443`)
	})

	t.Run("success", func(t *testing.T) {
		expectedTLSCert := "/etc/certs/cert.pem"
		expectedTLSKey := "/etc/certs/key.pem"
		expectedMTLSEnabled := "true"
		expectedServerAddr := ":443"
		expectedAllowedOrigins := []string{"https://arcade.arcadium.dev"}
		expectedAllowedMethods := []string{"GET", "OPTIONS"}
		expectedAllowedHeaders := []string{"content-type", "x-okta-user-agent-extended"}
		expectedReadTimeout := 7 * time.Second
		expectedWriteTimeout := 13 * time.Second
		expectedShutdownTimeout := 111 * time.Second

		t.Setenv("SERVER_ADDR", expectedServerAddr)
		t.Setenv("TLS_CERT", expectedTLSCert)
		t.Setenv("TLS_KEY", expectedTLSKey)
		t.Setenv("MTLS_ENABLED", expectedMTLSEnabled)
		t.Setenv("ALLOWED_ORIGINS", strings.Join(expectedAllowedOrigins, ", "))
		t.Setenv("ALLOWED_METHODS", strings.Join(expectedAllowedMethods, ", "))
		t.Setenv("ALLOWED_HEADERS", strings.Join(expectedAllowedHeaders, ", "))
		t.Setenv("READ_TIMEOUT", expectedReadTimeout.String())
		t.Setenv("WRITE_TIMEOUT", expectedWriteTimeout.String())
		t.Setenv("SHUTDOWN_TIMEOUT", expectedShutdownTimeout.String())

		cfg, err := server.NewConfig()

		assert.Nil(t, err)
		assert.Equal(t, cfg.TLSCert(), expectedTLSCert)
		assert.Equal(t, cfg.TLSKey(), expectedTLSKey)
		assert.Equal(t, cfg.ServerAddr(), expectedServerAddr)
		assert.Compare(t, cfg.AllowedOrigins(), expectedAllowedOrigins)
		assert.Compare(t, cfg.AllowedMethods(), expectedAllowedMethods)
		assert.Compare(t, cfg.AllowedHeaders(), expectedAllowedHeaders)
		assert.Equal(t, cfg.ReadTimeout(), expectedReadTimeout)
		assert.Equal(t, cfg.WriteTimeout(), expectedWriteTimeout)
		assert.Equal(t, cfg.ShutdownTimeout(), expectedShutdownTimeout)
	})
}
