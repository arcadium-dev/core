// Copyright 2021-2023 arcadium.dev <info@arcadium.dev>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package sql provides a database/sql wrapper.
package sql // import "arcadium.dev/core/sql"

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/rs/zerolog"
)

const (
	DefaultReconnectRetries = 5
	DefaultTxRetries        = 10
)

type (
	NullString = sql.NullString
	Result     = sql.Result
	Row        = sql.Row
	Rows       = sql.Rows
	Tx         = sql.Tx
	TxOptions  = sql.TxOptions
)

var (
	ErrNoRows = sql.ErrNoRows
)

type (
	// DB is a simple wrapper of sql.db.
	DB struct {
		driver           string
		connectURL       string
		db               *sql.DB
		middleware       []MiddlewareFunc
		reconnectRetries int
		txRetries        int
	}

	// Func is used to do some work.
	Func func(ctx context.Context) error

	// MiddlewareFunc is a function that receives a Func and returns a Func.
	MiddlewareFunc func(Func) Func
)

var (
	timeout time.Duration = 30 * time.Second

	// open allows for insertion of mock open functions.
	open = sql.Open

	// connect allows for insertion of mock connect functions.
	connect = func(ctx context.Context, db *sql.DB) error {
		count := 0

		retryCtx, retryCancel := context.WithTimeout(context.Background(), timeout)
		defer retryCancel()

		prevRetry, currRetry := time.Duration(1), time.Duration(1)
		for done := false; !done; {
			select {
			case <-time.After(currRetry * time.Second):
				nextRetry := currRetry + prevRetry
				prevRetry, currRetry = currRetry, nextRetry

				if err := db.PingContext(retryCtx); err != nil {
					count++
					zerolog.Ctx(ctx).Info().Int("count", count).Str("error", err.Error()).Msgf("ping failed, retrying...")
					continue
				}
				done = true

			case <-retryCtx.Done():
				return retryCtx.Err()
			}
		}

		return nil
	}
)

// Open opens a database specified by its database driver name and a
// driver-specific connect url.
func Open(ctx context.Context, driver, connectURL string, opts ...Option) (*DB, error) {
	db := &DB{
		driver:     driver,
		connectURL: connectURL,
	}

	o := &Options{
		ReconnectRetries: DefaultReconnectRetries,
		TxRetries:        DefaultTxRetries,
	}
	for _, opt := range opts {
		opt.Apply(o)
	}
	if o.ReconnectEnabled {
		db.middleware = append(db.middleware, db.reconnect)
		db.reconnectRetries = o.ReconnectRetries
	}
	db.txRetries = o.TxRetries

	sqlDB, err := open(db.driver, db.connectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := connect(ctx, sqlDB); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	db.db = sqlDB

	var (
		user, host string
	)

	u, err := url.Parse(connectURL)
	if err == nil {
		user = u.User.Username()
		host = u.Host
	}
	zerolog.Ctx(ctx).Info().Msgf("connected to database, driver '%s', user '%s', host '%s'", driver, user, host)

	return db, nil
}

// Ping checks DB status for the health checker
func (db *DB) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

// Close closes the database and prevents new queries from starting.
func (db *DB) Close() error {
	return db.db.Close()
}

// Exec executes a query without returning any rows.
func (db *DB) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	var result Result

	err := db.run(ctx, query, func(c context.Context) error {
		var e error
		result, e = db.db.ExecContext(c, query, args...)
		return e
	})

	return result, err
}

// Query executes a query that returns rows, typically a SELECT.
func (db *DB) Query(ctx context.Context, query string, args ...any) (*Rows, error) {
	var rows *Rows

	err := db.run(ctx, query, func(c context.Context) error {
		var e error
		rows, e = db.db.QueryContext(c, query, args...)
		return e
	})

	return rows, err
}

// QueryRow executes a query that is expected to return at most one row.
func (db *DB) QueryRow(ctx context.Context, query string, args ...any) *Row {
	var row *Row

	db.run(ctx, query, func(c context.Context) error {
		row = db.db.QueryRowContext(ctx, query, args...)
		return nil
	})

	return row
}

// BeginTx starts a transaction.
func (db *DB) BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error) {
	var tx *Tx

	err := db.run(ctx, "", func(c context.Context) error {
		var e error
		tx, e = db.db.BeginTx(ctx, opts)
		return e
	})

	return tx, err
}

// Tx wraps a transaction, implementing retries.
func (db *DB) Tx(ctx context.Context, opts *TxOptions, fn func(tx *Tx) error) error {
	maxRetries := db.txRetries
	retryCount := 0
	for {
		err := db.executeTx(ctx, opts, fn)
		if err == nil || !IsRetryable(err) {
			return err
		}

		retryCount++
		if retryCount >= maxRetries {
			return fmt.Errorf("transaction failed: exceeded retries: %w", err)
		}
	}
}

func (db *DB) executeTx(ctx context.Context, opts *TxOptions, fn func(tx *Tx) error) (err error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		r := recover()

		if r == nil && err == nil {
			err = tx.Commit()
			return
		}

		e := tx.Rollback()
		if e != nil && !errors.Is(e, sql.ErrTxDone) {
			zerolog.Ctx(ctx).Err(e).Msg("failed to rollback")
		}

		if r != nil {
			panic(r)
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	return nil
}

func (db *DB) run(ctx context.Context, query string, f Func) error {
	// Build a chain of functions, starting with the passed in function as the
	// last link in the chain. The goal is to have a single function that will
	// call each middleware function, ending with the passed in function.
	chain := f

	// Starting at the end of the slice and working backwards, link the functions
	// together.
	for i := len(db.middleware) - 1; i >= 0; i-- {
		chain = db.middleware[i](chain)
	}

	// Call the chain.
	return chain(ctx)
}

func (db *DB) reconnect(f Func) Func {
	return func(ctx context.Context) error {
		var err error
		logger := zerolog.Ctx(ctx)

		// This implements the connection retry loop as documented here:
		// https://www.cockroachlabs.com/docs/stable/node-shutdown.html#connection-retry-loop

		maxRetries := db.reconnectRetries
		retryCount := 0
		for {
			err = f(ctx)

			retryCount++
			if retryCount > maxRetries {
				logger.Error().Msg("db reconnect, exceeding retry count")
				return err
			}

			if !IsAdminShutdown(err) {
				break
			}

			logger.Info().Msg("database admin shutdown detected, reconnecting...")

			// Open a new connection. This will be routed to the non-draining node.
			sqlDB, e := open(db.driver, db.connectURL)
			if e != nil {
				logger.Warn().Msgf("failed to open new db after admin shutdown, reason: %s", e)
				continue
			}
			if e := connect(ctx, sqlDB); e != nil {
				sqlDB.Close()
				logger.Warn().Msgf("failed to connect to new db after admin shutdown, reason: %s", e)
				continue
			}
			logger.Info().Msg("database admin shutdown detected, reconnected")

			// Close the prior connection after we have a known good connection.
			if e := db.db.Close(); e != nil {
				logger.Warn().Msgf("failed to close existing db connection after admin shutdown, reason: %s", e)
			}
			db.db = sqlDB
			err = nil
		}
		return err
	}
}

func IsAdminShutdown(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == "unexpected EOF" {
		return true
	}
	code := ErrCode(err)
	return code == "57P01"
}

func IsRetryable(err error) bool {
	// Check for standard PG errcode SerializationFailureError:40001.
	return ErrCode(err) == "40001"
}

func ErrCode(err error) string {
	if e, ok := err.(interface{ SQLState() string }); ok {
		return e.SQLState()
	}
	return ""
}
