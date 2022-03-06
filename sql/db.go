// Copyright 2021-2022 arcadium.dev <info@arcadium.dev>
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

package sql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"
)

// Open opens a database specified by its database driver name and a
// driver-specific connect url.
func Open(driver, connectURL string, logger Logger) (*DB, error) {
	db, err := open(driver, connectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s database: %w", driver, err)
	}

	if err := connect(db, logger); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	var (
		user, host string
	)

	u, err := url.Parse(connectURL)
	if err == nil {
		user = u.User.Username()
		host = u.Host
	}
	logger.Info("msg", "connected to database", "driver", driver, "user", user, "host", host)

	return &DB{DB: db}, nil
}

type (
	// DB is a simple wrapper of sql.DB.
	DB struct {
		*sql.DB
	}

	// Logger defines the logger needed by the sql package.
	Logger interface {
		Info(...interface{})
	}
)

var (
	timeout time.Duration = 30 * time.Second

	// open allows for insertion of mock open functions.
	open = sql.Open

	// connect allows for insertion of mock connect functions.
	connect = func(db *sql.DB, logger Logger) error {
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
					logger.Info("msg", "ping failed, retrying...", "count", count, "error", err.Error())
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
