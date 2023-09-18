package log

import (
	"testing"

	"github.com/rs/zerolog"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/require"
)

func TestAsDefault(t *testing.T) {
	var opts options

	o := AsDefault()
	require.NotNil(t, o)

	o.apply(&opts)
	assert.True(t, opts.asDefault)
}

func TestWithLevel(t *testing.T) {
	levelCheck := func(t *testing.T, lvl zerolog.Level) {
		var opts options

		o := WithLevel(lvl)
		require.NotNil(t, o)

		o.apply(&opts)
		assert.Equal(t, opts.level, lvl)
	}

	for _, lvl := range []zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.WarnLevel, zerolog.ErrorLevel} {
		levelCheck(t, lvl)
	}
}

func TestWithLevelFieldName(t *testing.T) {
	var opts options
	name := "severity"

	o := WithLevelFieldName(name)
	require.NotNil(t, o)

	o.apply(&opts)
	assert.Equal(t, name, opts.levelFieldName)
}

func TestWithMessageFieldName(t *testing.T) {
	var opts options
	name := "msg"

	o := WithMessageFieldName(name)
	require.NotNil(t, o)

	o.apply(&opts)
	assert.Equal(t, name, opts.messageFieldName)
}

func TestWithOutput(t *testing.T) {
	var (
		opts options
		b    = NewStringBuffer()
	)
	o := WithOutput(b)
	require.NotNil(t, o)

	o.apply(&opts)
	assert.Equal(t, opts.writer.(*StringBuffer), b)
}

func TestWithoutTimestamp(t *testing.T) {
	var opts options

	o := WithoutTimestamp()
	require.NotNil(t, o)

	o.apply(&opts)
	assert.Equal(t, opts.timestamped, false)
}

func TestWithTimestampFieldName(t *testing.T) {
	var opts options
	name := "ts"

	o := WithTimestampFieldName(name)
	require.NotNil(t, o)

	o.apply(&opts)
	assert.Equal(t, name, opts.timestampFieldName)
}
