package sql_test

import (
	"testing"

	"arcadium.dev/core/sql"
)

func TestWithReconnect(t *testing.T) {
	t.Run("with invalid retries", func(t *testing.T) {
		o := &sql.Options{
			ReconnectRetries: sql.DefaultReconnectRetries,
		}
		sql.WithReconnect(-1).Apply(o)

		if o.ReconnectEnabled != true {
			t.Errorf("Reconnect was not enabled")
		}
		if o.ReconnectRetries != sql.DefaultReconnectRetries {
			t.Errorf("Unexpected reconnect retries: %d", o.ReconnectRetries)
		}
	})

	t.Run("with valid retries", func(t *testing.T) {
		o := &sql.Options{
			ReconnectRetries: sql.DefaultReconnectRetries,
		}
		sql.WithReconnect(42).Apply(o)

		if o.ReconnectEnabled != true {
			t.Errorf("Reconnect was not enabled")
		}
		if o.ReconnectRetries != 42 {
			t.Errorf("Unexpected reconnect retries: %d", o.ReconnectRetries)
		}
	})
}

func TestWithTxRetries(t *testing.T) {
	t.Run("with invalid retries", func(t *testing.T) {
		o := &sql.Options{
			TxRetries: sql.DefaultTxRetries,
		}
		sql.WithTxRetries(-1).Apply(o)

		if o.TxRetries != sql.DefaultTxRetries {
			t.Errorf("Unexpected tx retries: %d", o.TxRetries)
		}
	})

	t.Run("with valid retries", func(t *testing.T) {
		o := &sql.Options{
			TxRetries: sql.DefaultTxRetries,
		}
		sql.WithTxRetries(42).Apply(o)

		if o.TxRetries != 42 {
			t.Errorf("Unexpected tx retries: %d", o.TxRetries)
		}
	})
}
