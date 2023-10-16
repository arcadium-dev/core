package sql

import (
	"context"
	"database/sql"
	"time"
)

var (
// ContextWithQuery = contextWithQuery
// QueryFromContext = queryFromContext
)

type OpenFunc func(driverName, dataSourceName string) (*sql.DB, error)

func SetOpen(newOpen OpenFunc) OpenFunc {
	oldOpen := open
	open = newOpen
	return oldOpen
}

type ConnectFunc func(ctx context.Context, db *sql.DB) error

func SetConnect(newConnect ConnectFunc) ConnectFunc {
	oldConnect := connect
	connect = newConnect
	return oldConnect
}

func SetTimeout(newTimeout time.Duration) time.Duration {
	oldTimeout := timeout
	timeout = newTimeout
	return oldTimeout
}

func mockMiddleware(f Func) Func {
	return func(ctx context.Context) error { return f(ctx) }
}

func NewDB(db *sql.DB, retries ...int) *DB {
	reconnectRetries := DefaultReconnectRetries
	if len(retries) >= 1 && retries[0] > 0 {
		reconnectRetries = retries[0]
	}

	txRetries := DefaultTxRetries
	if len(retries) >= 2 && retries[1] > 0 {
		txRetries = retries[1]
	}

	return &DB{
		db:               db,
		middleware:       []MiddlewareFunc{mockMiddleware, mockMiddleware},
		reconnectRetries: reconnectRetries,
		txRetries:        txRetries,
	}
}

func (db *DB) SetDB(sqlDB *sql.DB) {
	db.db = sqlDB
}

func (db *DB) Reconnect(f Func) Func {
	return db.reconnect(f)
}
