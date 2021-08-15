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

import (
	"context"
	"database/sql"
	"time"

	"arcadium.dev/core/errors"
)

// Open opens a database specified by its config. The config will
// provide database driver name and a driver-specific data source name.
func Open(config Config, opts ...Option) (*sql.DB, error) {
	db, err := sql.Open(config.DriverName(), config.DSN())
	if err != nil {
		return nil, err
	}

	// Set options
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}

	// If there is a migration, run it against the sql db.
	if o.migration != nil {
		if err := Connect(db); err != nil {
			return nil, err
		}
		if err := o.migration(db); err != nil {
			return nil, err
		}
	}
	return db, nil
}

const (
	timeout = 30
)

// Connect establishes a connection to the database. It will
func Connect(db *sql.DB) error {
	retryCtx, retryCancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer retryCancel()

	prevRetry, currRetry := time.Duration(1), time.Duration(1)

	for {
		select {
		case <-time.After(currRetry * time.Second):
			nextRetry := currRetry + prevRetry
			prevRetry, currRetry = currRetry, nextRetry

			if err := db.PingContext(retryCtx); err != nil {
				continue
			}
			return nil

		case <-retryCtx.Done():
			return errors.Wrapf(retryCtx.Err(), "failed to connect to the database")
		}
	}
}
