// Copyright 2021 Ian Cahoon <icahoon@gmail.com>
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

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

// sqlDB

type sqlDB struct {
	*sql.DB
}

func (db *sqlDB) Begin(ctx context.Context, opts *TxOptions) (Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &sqlTx{Tx: tx}, nil
}

func (db *sqlDB) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	result, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (db *sqlDB) Prepare(ctx context.Context, query string) (Stmt, error) {
	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &sqlStmt{Stmt: stmt}, nil
}

func (db *sqlDB) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return rows, err
}

func (db *sqlDB) QueryRow(ctx context.Context, query string, args ...interface{}) *Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}

func (db *sqlDB) Close() error {
	return errors.WithStack(db.DB.Close())
}

func (db *sqlDB) Ping(ctx context.Context) error {
	return errors.WithStack(db.DB.PingContext(ctx))
}

// sqlStmt

type sqlStmt struct {
	*sql.Stmt
}

func (stmt *sqlStmt) Exec(ctx context.Context, args ...interface{}) (Result, error) {
	result, err := stmt.Stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (stmt *sqlStmt) Query(ctx context.Context, args ...interface{}) (*Rows, error) {
	rows, err := stmt.Stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return rows, err
}

func (stmt *sqlStmt) QueryRow(ctx context.Context, args ...interface{}) *Row {
	return stmt.Stmt.QueryRowContext(ctx, args...)
}

func (stmt *sqlStmt) Close() error {
	return errors.WithStack(stmt.Stmt.Close())
}

// sqlTx

type sqlTx struct {
	*sql.Tx
}

func (tx *sqlTx) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	result, err := tx.Tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (tx *sqlTx) Prepare(ctx context.Context, query string) (Stmt, error) {
	stmt, err := tx.Tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &sqlStmt{Stmt: stmt}, nil
}

func (tx *sqlTx) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, err := tx.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return rows, err
}

func (tx *sqlTx) QueryRow(ctx context.Context, query string, args ...interface{}) *Row {
	return tx.Tx.QueryRowContext(ctx, query, args...)
}

func (tx *sqlTx) Commit() error {
	return errors.WithStack(tx.Tx.Commit())
}

func (tx *sqlTx) Rollback() error {
	return errors.WithStack(tx.Tx.Rollback())
}

func (tx *sqlTx) Stmt(ctx context.Context, stmt Stmt) Stmt {
	return &sqlStmt{Stmt: tx.Tx.StmtContext(ctx, stmt.(*sqlStmt).Stmt)}
}
