package sql

import (
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOpen(t *testing.T) {
	sqlmock.MonitorPingsOption(true)
	logger := mockLogger{}

	t.Run("open failure", func(t *testing.T) {
		// Replace open with a mock that returns an error.
		open = func(driver, url string) (*sql.DB, error) {
			return nil, errors.New("open failure")
		}

		_, err := Open("postgres", "url", logger)
		if err == nil || err.Error() != "failed to open postgres database: open failure" {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("connect failure", func(t *testing.T) {
		// Replace open with a mock that returns an sqlmock db.
		open = func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock")
			}
			return db, nil
		}

		// Replace connect with a mock that will return an error.
		oldConnect := connect
		connect = func(*sql.DB, Logger) error {
			return errors.New("connect failure")
		}
		defer func() { connect = oldConnect }()

		_, err := Open("driver", "url", logger)
		if err == nil || err.Error() != "failed to connect to the database: connect failure" {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("connect failure - timeout", func(t *testing.T) {
		open = func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			if db == nil || err != nil {
				t.Fatal("failed to create sqlmock")
			}
			return db, nil
		}

		// Timeout is less than the retry, so we will hit the context deadline.
		timeout = 50 * time.Millisecond
		db, err := Open("driver", "url", logger)

		if db != nil {
			t.Errorf("Unexpected db: %+v", db)
		}
		if err == nil {
			t.Error("Expected an error")
		}
		if strings.Contains(err.Error(), "failed to connect to database") {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("connect success", func(t *testing.T) {
		open = func(driver, url string) (*sql.DB, error) {
			db, _, err := sqlmock.New()
			if db == nil || err != nil {
				t.Fatal("failed to create sqlmock")
			}
			return db, nil
		}

		// Timeout is greater than the retry, so we will get a successful ping.
		timeout = 2 * time.Second
		db, err := Open("driver", "url", logger)

		if db == nil {
			t.Error("Expected a db")
		}
		db.Close()
		if err != nil {
			t.Error("Unexpected error")
		}
	})
}

type (
	mockLogger struct{}
)

func (m mockLogger) Info(...interface{}) {}
