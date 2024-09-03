package mpserver_test

import (
	"os"
	"testing"
	"time"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/mpserver"
)

func TestNewConfig(t *testing.T) {
	for _, env := range []string{"DSN", "LOG_LEVEL"} {
		os.Unsetenv(env)
	}

	t.Run("test defaults", func(t *testing.T) {
		cfg, err := mpserver.NewConfig()

		assert.Nil(t, err)
		assert.Equal(t, cfg.LogLevel(), mpserver.DefaultLogLevel)
		assert.Equal(t, cfg.ShutdownTimeout(), mpserver.DefaultShutdownTimeout)
	})

	t.Run("success", func(t *testing.T) {
		expectedLogLevel := "warn"
		expectedShutdownTimeout := 111 * time.Second

		t.Setenv("LOG_LEVEL", expectedLogLevel)
		t.Setenv("SHUTDOWN_TIMEOUT", "111s")

		cfg, err := mpserver.NewConfig()

		assert.Nil(t, err)
		assert.Equal(t, cfg.LogLevel(), expectedLogLevel)
		assert.Equal(t, cfg.ShutdownTimeout(), expectedShutdownTimeout)
	})
}
