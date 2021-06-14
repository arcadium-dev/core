// Copyright 2021 arcadium.dev <info@arcadium.dev>
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

package sql // import "arcadium.dev/core/sql"

//go:generate mockgen -package mocksql -destination ./mock/sql.go . DB,Stmt,Tx

import (
	"context"
	"database/sql"

	"arcadium.dev/core/log"
)

type Result = sql.Result
type Row = sql.Row
type Rows = sql.Rows
type TxOptions = sql.TxOptions

type DB interface {
	Begin(ctx context.Context, opts *TxOptions) (Tx, error)

	Exec(ctx context.Context, query string, args ...interface{}) (Result, error)
	Prepare(ctx context.Context, query string) (Stmt, error)
	Query(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *Row

	Close() error
	Ping(ctx context.Context) error
}

type Stmt interface {
	Exec(ctx context.Context, args ...interface{}) (Result, error)
	Query(ctx context.Context, args ...interface{}) (*Rows, error)
	QueryRow(ctx context.Context, args ...interface{}) *Row

	Close() error
}

type Tx interface {
	Exec(ctx context.Context, query string, args ...interface{}) (Result, error)
	Prepare(ctx context.Context, query string) (Stmt, error)
	Query(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *Row

	Commit() error
	Rollback() error

	Stmt(ctx context.Context, stmt Stmt) Stmt
}

/*
	What I would like:

	type Logger[L any] interface {
		WithField(key string, value interface{}) L

		Infoln(args ...interface{})
		Infof(format string, args ...interface{})
	}
*/
type Logger = log.Logger
