package rest_test

import (
	"os"
	"strings"
	"testing"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/rest"
)

func TestNewConfig(t *testing.T) {
	for _, env := range []string{"DSN", "LOG_LEVEL"} {
		os.Unsetenv(env)
	}

	t.Run("test defaults", func(t *testing.T) {
		t.Setenv("TEST_SERVER_ADDR", ":8443")

		cfg, err := rest.NewConfig("test")

		assert.Nil(t, err)
		assert.Equal(t, cfg.LogLevel(), rest.DefaultLogLevel)
	})

	t.Run("success", func(t *testing.T) {
		expectedLogLevel := "warn"
		expectedTLSCert := "/etc/certs/cert.pem"
		expectedTLSKey := "/etc/certs/key.pem"
		expectedMTLSEnabled := "true"
		expectedServerAddr := ":443"
		expectedAllowedOrigins := []string{"https://arcade.arcadium.dev"}
		expectedAllowedMethods := []string{"GET", "OPTIONS"}
		expectedAllowedHeaders := []string{"content-type", "x-okta-user-agent-extended"}
		expectedPProfEnabled := "true"

		t.Setenv("LOG_LEVEL", expectedLogLevel)
		t.Setenv("TLS_CERT", expectedTLSCert)
		t.Setenv("TLS_KEY", expectedTLSKey)
		t.Setenv("MTLS_ENABLED", expectedMTLSEnabled)
		t.Setenv("SERVER_ADDR", expectedServerAddr)
		t.Setenv("ALLOWED_ORIGINS", strings.Join(expectedAllowedOrigins, ", "))
		t.Setenv("ALLOWED_METHODS", strings.Join(expectedAllowedMethods, ", "))
		t.Setenv("ALLOWED_HEADERS", strings.Join(expectedAllowedHeaders, ", "))
		t.Setenv("PPROF_ENABLED", expectedPProfEnabled)

		cfg, err := rest.NewConfig()

		assert.Nil(t, err)
		assert.Equal(t, cfg.LogLevel(), expectedLogLevel)
		assert.Equal(t, cfg.TLSCert(), expectedTLSCert)
		assert.Equal(t, cfg.TLSKey(), expectedTLSKey)
		assert.Equal(t, cfg.ServerAddr(), expectedServerAddr)
		assert.Compare(t, cfg.AllowedOrigins(), expectedAllowedOrigins)
		assert.Compare(t, cfg.AllowedMethods(), expectedAllowedMethods)
		assert.Compare(t, cfg.AllowedHeaders(), expectedAllowedHeaders)
		assert.True(t, cfg.PProfEnabled())
	})
}
