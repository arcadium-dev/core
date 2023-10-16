package sql_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/log"
	"arcadium.dev/core/require"
	csql "arcadium.dev/core/sql"
)

func TestOpen(t *testing.T) {
	ctx := context.Background()

	sqlmock.MonitorPingsOption(true)

	t.Run("open failure", func(t *testing.T) {
		// Replace open with a mock that returns an error.
		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			return nil, errors.New("open failure")
		})
		defer func() { csql.SetOpen(origOpen) }()

		_, err := csql.Open(ctx, "driver", "url")
		assert.Error(t, err, "failed to open database: open failure")
	})

	t.Run("connect failure", func(t *testing.T) {
		// Replace open with a mock that returns an sqlmock db.
		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			require.Nil(t, err)
			return db, nil
		})
		defer func() { csql.SetOpen(origOpen) }()

		// Replace connect with a mock that will return an error.
		origConnect := csql.SetConnect(func(ctx context.Context, db *sql.DB) error {
			return errors.New("connect failure")
		})
		defer func() { csql.SetConnect(origConnect) }()

		_, err := csql.Open(ctx, "driver", "url")
		assert.Error(t, err, "failed to connect to the database: connect failure")
	})

	t.Run("connect failure - timeout", func(t *testing.T) {
		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			require.Nil(t, err)
			return db, nil
		})
		defer func() { csql.SetOpen(origOpen) }()

		// Timeout is less than the retry, so we will hit the context deadline.
		origTimeout := csql.SetTimeout(50 * time.Millisecond)
		defer func() { csql.SetTimeout(origTimeout) }()

		_, err := csql.Open(ctx, "driver", "url")

		assert.Error(t, err, "failed to connect to the database: context deadline exceeded")
	})

	t.Run("connect success", func(t *testing.T) {
		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			assert.Nil(t, err)
			return db, nil
		})
		defer func() { csql.SetOpen(origOpen) }()

		// Timeout is greater than the retry, so we will get a successful ping.
		origTimeout := csql.SetTimeout(2 * time.Second)
		defer func() { csql.SetTimeout(origTimeout) }()

		db, err := csql.Open(ctx, "driver", "url")

		assert.Nil(t, err)
		db.Close()
	})

	t.Run("with reconnect option", func(t *testing.T) {
		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			require.Nil(t, err)
			return db, nil
		})
		defer func() { csql.SetOpen(origOpen) }()

		// Timeout is greater than the retry, so we will get a successful ping.
		origTimeout := csql.SetTimeout(2 * time.Second)
		defer func() { csql.SetTimeout(origTimeout) }()

		db, err := csql.Open(ctx, "driver", "url", csql.WithReconnect(42))

		assert.Nil(t, err)
		db.Close()
	})

	t.Run("with tx retries", func(t *testing.T) {
		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			require.Nil(t, err)
			return db, nil
		})
		defer func() { csql.SetOpen(origOpen) }()

		// Timeout is greater than the retry, so we will get a successful ping.
		origTimeout := csql.SetTimeout(2 * time.Second)
		defer func() { csql.SetTimeout(origTimeout) }()

		db, err := csql.Open(ctx, "driver", "url", csql.WithTxRetries(42))

		assert.Nil(t, err)
		db.Close()
	})
}

func TestDBPing(t *testing.T) {
	ctx := context.Background()

	t.Run("ping failure", func(t *testing.T) {
		mdb, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectPing().WillReturnError(errors.New("ping failure"))
		db := csql.NewDB(mdb)

		err = db.Ping(ctx)

		assert.Error(t, err, `ping failure`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("success", func(t *testing.T) {
		mdb, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectPing()
		db := csql.NewDB(mdb)

		err = db.Ping(context.Background())

		assert.Nil(t, err)
		assert.MockExpectationsMet(t, mock)
	})
}

func TestDBClose(t *testing.T) {
	t.Run("close failure", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		mock.ExpectClose().WillReturnError(errors.New("close failure"))
		db := csql.NewDB(mdb)

		err = db.Close()

		assert.Error(t, err, `close failure`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("success", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		mock.ExpectClose()
		db := csql.NewDB(mdb)

		err = db.Close()

		assert.Nil(t, err)
		assert.MockExpectationsMet(t, mock)
	})
}

func TestDBExec(t *testing.T) {
	ctx := context.Background()

	t.Run("exec failure", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectExec(`^INSERT INTO foo (.+) VALUES (.+)$`).WillReturnError(errors.New("exec failure"))
		db := csql.NewDB(mdb)

		_, err = db.Exec(ctx, "INSERT INTO foo (bar) VALUES ($1)", "baz")

		assert.Error(t, err, `exec failure`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("success", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectExec(`^INSERT INTO foo (.+) VALUES (.+)$`).WillReturnResult(sqlmock.NewResult(0, 1))
		db := csql.NewDB(mdb)

		_, err = db.Exec(ctx, "INSERT INTO foo (bar) VALUES ($1)", "baz")

		assert.Nil(t, err)
		assert.MockExpectationsMet(t, mock)
	})
}

func TestDBQuery(t *testing.T) {
	ctx := context.Background()

	t.Run("query failure", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectQuery(`^SELECT foo FROM bar$`).WillReturnError(errors.New("query failure"))
		db := csql.NewDB(mdb)

		_, err = db.Query(ctx, "SELECT foo FROM bar")

		assert.Error(t, err, `query failure`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("success", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		rows := sqlmock.NewRows([]string{"foo"}).AddRow("baz")
		mock.ExpectQuery(`^SELECT foo FROM bar$`).WillReturnRows(rows)
		db := csql.NewDB(mdb)

		_, err = db.Query(ctx, "SELECT foo FROM bar")

		assert.Nil(t, err)
		assert.MockExpectationsMet(t, mock)
	})
}

func TestDBQueryRow(t *testing.T) {
	t.Run("queryrow failure", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectQuery(`^SELECT foo FROM bar$`).WillReturnError(errors.New("queryrow failure"))
		db := csql.NewDB(mdb)

		row := db.QueryRow(context.Background(), "SELECT foo FROM bar")
		var value string
		err = row.Scan(&value)

		assert.Error(t, err, `queryrow failure`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("success", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		rows := sqlmock.NewRows([]string{"foo"}).AddRow("baz")
		mock.ExpectQuery(`^SELECT foo FROM bar$`).WillReturnRows(rows)
		db := csql.NewDB(mdb)

		row := db.QueryRow(context.Background(), "SELECT foo FROM bar")
		var value string
		err = row.Scan(&value)

		assert.Nil(t, err)
		assert.Equal(t, value, "baz")
		assert.MockExpectationsMet(t, mock)
	})
}

func TestDBBeginTx(t *testing.T) {
	ctx := context.Background()

	t.Run("begintx failure", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectBegin().WillReturnError(errors.New("begintx failure"))
		db := csql.NewDB(mdb)

		_, err = db.BeginTx(ctx, nil)

		assert.Error(t, err, `begintx failure`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("success", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectBegin()
		db := csql.NewDB(mdb)

		_, err = db.BeginTx(context.Background(), nil)

		assert.Nil(t, err)
		assert.MockExpectationsMet(t, mock)
	})
}

func TestDBTx(t *testing.T) {
	ctx := context.Background()

	t.Run("begintx failure", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectBegin().WillReturnError(errors.New("begintx failure"))
		db := csql.NewDB(mdb)

		err = db.Tx(ctx, nil, func(tx *csql.Tx) error {
			return nil
		})

		assert.Error(t, err, `begintx failure`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("exceed retries failure", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		for i := 0; i < csql.DefaultTxRetries; i++ {
			mock.ExpectBegin()
		}
		db := csql.NewDB(mdb)

		err = db.Tx(context.Background(), nil, func(tx *csql.Tx) error {
			return mockPgError{code: "40001"}
		})

		assert.Error(t, err, `transaction failed: exceeded retries: SQLSTATE 40001`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("fn panic - rollback", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectBegin()
		mock.ExpectRollback()
		db := csql.NewDB(mdb)

		defer func() {
			r := recover()

			actual, ok := r.(string)
			assert.True(t, ok)
			assert.Equal(t, actual, "woohoo!")
			assert.MockExpectationsMet(t, mock)
		}()

		db.Tx(ctx, nil, func(tx *csql.Tx) error {
			panic("woohoo!")
		})
	})

	t.Run("fn failure - rollback error", func(t *testing.T) {
		ctx, b := log.SetupTestLogging(t)

		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectBegin()
		mock.ExpectRollback().WillReturnError(errors.New("rollback failure"))
		db := csql.NewDB(mdb)

		err = db.Tx(ctx, nil, func(tx *csql.Tx) error {
			return errors.New("fn failure")
		})

		assert.Error(t, err, `fn failure`)
		require.Equal(t, b.Len(), 1)
		assert.Equal(t, b.Index(0), `{"severity":"error","error":"rollback failure","message":"failed to rollback"}`+"\n")
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("fn failure - rollback success", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectBegin()
		mock.ExpectRollback()
		db := csql.NewDB(mdb)

		err = db.Tx(ctx, nil, func(tx *csql.Tx) error {
			return errors.New("fn failure")
		})

		assert.Error(t, err, `fn failure`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("fn success - commit failure", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(errors.New("commit failure"))
		db := csql.NewDB(mdb)

		err = db.Tx(context.Background(), nil, func(tx *csql.Tx) error {
			return nil
		})

		assert.Error(t, err, `commit failure`)
		assert.MockExpectationsMet(t, mock)
	})

	t.Run("fn success - commit success", func(t *testing.T) {
		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer mdb.Close()

		mock.ExpectBegin()
		mock.ExpectCommit()
		db := csql.NewDB(mdb)

		err = db.Tx(context.Background(), nil, func(tx *csql.Tx) error {
			return nil
		})

		assert.Nil(t, err)
		assert.MockExpectationsMet(t, mock)
	})
}

func TestDBReconnect(t *testing.T) {
	ctx := context.Background()

	t.Run("non admin shutdown error", func(t *testing.T) {
		// Replace open with a mock that returns an sqlmock db.
		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			require.Nil(t, err)
			return db, nil
		})
		defer func() { csql.SetOpen(origOpen) }()

		// Always connect.
		origConnect := csql.SetConnect(func(ctx context.Context, db *sql.DB) error {
			return nil
		})
		defer func() { csql.SetConnect(origConnect) }()

		db := csql.NewDB(nil, 50000)

		// Always fail with an admin shutdown.
		f := func(ctx context.Context) error {
			return errors.New("failure")
		}

		err := db.Reconnect(f)(ctx)

		assert.Error(t, err, "failure")
	})

	t.Run("open failure", func(t *testing.T) {
		ctx, b := log.SetupTestLogging(t)

		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			return nil, errors.New("always fail")
		})
		defer func() { csql.SetOpen(origOpen) }()

		db := csql.NewDB(nil, 2)

		// Always fail with an admin shutdown.
		f := func(ctx context.Context) error {
			return mockPgError{code: "57P01"}
		}

		err := db.Reconnect(f)(ctx)

		assert.Error(t, err, "SQLSTATE 57P01")
		require.Equal(t, b.Len(), 5)
		assert.Equal(t, b.Index(3), `{"severity":"warn","message":"failed to open new db after admin shutdown, reason: always fail"}`+"\n")
		assert.Equal(t, b.Index(4), `{"severity":"error","message":"db reconnect, exceeding retry count"}`+"\n")
	})

	t.Run("connect failure", func(t *testing.T) {
		ctx, b := log.SetupTestLogging(t)

		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			require.Nil(t, err)
			return db, nil
		})
		defer func() { csql.SetOpen(origOpen) }()

		origConnect := csql.SetConnect(func(ctx context.Context, db *sql.DB) error {
			return errors.New("always fail")
		})
		defer func() { csql.SetConnect(origConnect) }()

		db := csql.NewDB(nil, 2)

		// Always fail with an admin shutdown.
		f := func(ctx context.Context) error {
			return mockPgError{code: "57P01"}
		}

		err := db.Reconnect(f)(ctx)

		assert.Error(t, err, "SQLSTATE 57P01")
		require.Equal(t, b.Len(), 5)
		assert.Equal(t, b.Index(3), `{"severity":"warn","message":"failed to connect to new db after admin shutdown, reason: always fail"}`+"\n")
		assert.Equal(t, b.Index(4), `{"severity":"error","message":"db reconnect, exceeding retry count"}`+"\n")
	})

	t.Run("close failure", func(t *testing.T) {
		ctx, b := log.SetupTestLogging(t)

		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			require.Nil(t, err)
			return db, nil
		})
		defer func() { csql.SetOpen(origOpen) }()

		origConnect := csql.SetConnect(func(ctx context.Context, db *sql.DB) error {
			return nil
		})
		defer func() { csql.SetConnect(origConnect) }()

		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		mock.ExpectClose().WillReturnError(errors.New("close failure"))

		db := csql.NewDB(nil, 1)
		db.SetDB(mdb)

		// Always fail with an admin shutdown.
		f := func(ctx context.Context) error {
			return mockPgError{code: "57P01"}
		}

		err = db.Reconnect(f)(ctx)

		assert.Error(t, err, "SQLSTATE 57P01")
		assert.MockExpectationsMet(t, mock)
		require.Equal(t, b.Len(), 4)
		assert.Equal(t, b.Index(2), `{"severity":"warn","message":"failed to close existing db connection after admin shutdown, reason: close failure"}`+"\n")
		assert.Equal(t, b.Index(3), `{"severity":"error","message":"db reconnect, exceeding retry count"}`+"\n")
	})

	t.Run("success", func(t *testing.T) {
		ctx, b := log.SetupTestLogging(t)

		origOpen := csql.SetOpen(func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			assert.Nil(t, err)
			return db, nil
		})
		defer func() { csql.SetOpen(origOpen) }()

		origConnect := csql.SetConnect(func(ctx context.Context, db *sql.DB) error {
			return nil
		})
		defer func() { csql.SetConnect(origConnect) }()

		mdb, mock, err := sqlmock.New()
		require.Nil(t, err)
		mock.ExpectClose()

		db := csql.NewDB(nil, 5)
		db.SetDB(mdb)

		i := 0
		f := func(ctx context.Context) error {
			i++
			if i <= 1 {
				return mockPgError{code: "57P01"}
			}
			return nil
		}

		err = db.Reconnect(f)(ctx)
		assert.Nil(t, err)

		err = db.Reconnect(f)(ctx)
		assert.Nil(t, err)

		assert.MockExpectationsMet(t, mock)
		require.Equal(t, b.Len(), 2)
		assert.Equal(t, b.Index(0), `{"severity":"info","message":"database admin shutdown detected, reconnecting..."}`+"\n")
		assert.Equal(t, b.Index(1), `{"severity":"info","message":"database admin shutdown detected, reconnected"}`+"\n")
	})
}

func TestIsAdminShutdown(t *testing.T) {
	assert.False(t, csql.IsAdminShutdown(nil))
	assert.True(t, csql.IsAdminShutdown(mockPgError{code: "57P01"}))
	assert.True(t, csql.IsAdminShutdown(errors.New("unexpected EOF")))
	assert.False(t, csql.IsAdminShutdown(errors.New("foobar")))
}

type (
	mockPgError struct {
		code string
	}
)

func (m mockPgError) SQLState() string {
	return m.code
}

func (m mockPgError) Error() string {
	return "SQLSTATE " + m.code
}
