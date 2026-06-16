package telnet_test

import (
	"testing"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/telnet"
)

var (
	expectedServerAddr = ":2323"
)

func setenv(t *testing.T) {
	t.Setenv("SERVER_ADDR", expectedServerAddr)
}

func TestServerConfig_New(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setenv(t)
		cfg, err := telnet.NewServerConfig()

		assert.Nil(t, err)
		assert.Equal(t, cfg.ServerAddr(), expectedServerAddr)
	})
}

func TestServerConfig_ToOptions(t *testing.T) {
	setenv(t)
	cfg, err := telnet.NewServerConfig()
	assert.Nil(t, err)
	options := cfg.ToOptions()

	assert.Equal(t, len(options), 1)
}
