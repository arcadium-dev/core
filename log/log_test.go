package log_test

import (
	"os"
	"testing"

	"github.com/rs/zerolog"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/log"
)

func TestNew(t *testing.T) {
	t.Run("Test invalid level option", func(t *testing.T) {
		_, err := log.New(log.WithLevel(zerolog.Level(42)))
		assert.IsError(t, err, log.ErrInvalidLevel)
	})

	t.Run("Test invalid output option", func(t *testing.T) {
		_, err := log.New(log.WithOutput(nil))
		assert.IsError(t, err, log.ErrInvalidOutput)
	})

	t.Run("Test valid level option", func(t *testing.T) {
		l, err := log.New(log.WithLevel(zerolog.ErrorLevel))
		assert.Nil(t, err)
		assert.NotNil(t, l)
	})

	t.Run("Test valid output option", func(t *testing.T) {
		l, err := log.New(log.WithOutput(log.NewStringBuffer()))
		assert.Nil(t, err)
		assert.NotNil(t, l)
	})

	t.Run("Test all options", func(t *testing.T) {
		levelName := "severity"
		messageName := "msg"
		timestampName := "ts"

		l, err := log.New(
			log.AsDefault(),
			log.WithLevel(zerolog.WarnLevel),
			log.WithLevelFieldName(levelName),
			log.WithMessageFieldName(messageName),
			log.WithOutput(os.Stdout),
			log.WithoutTimestamp(),
			log.WithTimestampFieldName(timestampName),
		)

		assert.Nil(t, err)
		assert.NotNil(t, l)

		assert.Equal(t, levelName, zerolog.LevelFieldName)
		assert.Equal(t, messageName, zerolog.MessageFieldName)
		assert.Equal(t, timestampName, zerolog.TimestampFieldName)
	})
}

func TestToLevel(t *testing.T) {
	levels := []struct {
		s string
		l zerolog.Level
	}{
		{s: "", l: zerolog.InfoLevel},
		{s: "info", l: zerolog.InfoLevel},
		{s: "debug", l: zerolog.DebugLevel},
		{s: "DEBUG", l: zerolog.DebugLevel},
		{s: "INFO", l: zerolog.InfoLevel},
		{s: "warn", l: zerolog.WarnLevel},
		{s: "WARN", l: zerolog.WarnLevel},
		{s: "error", l: zerolog.ErrorLevel},
		{s: "ERROR", l: zerolog.ErrorLevel},
		{s: "invalid", l: zerolog.NoLevel},
	}
	for _, l := range levels {
		assert.Equal(t, log.ToLevel(l.s), l.l)
	}
}
