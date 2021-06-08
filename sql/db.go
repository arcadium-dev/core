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
	"database/sql"

	"github.com/pkg/errors"
)

func Open(config Config, opts ...Option) (DB, error) {
	sqldb, err := sql.Open(config.DriverName(), config.DSN())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var db DB = &sqlDB{DB: sqldb}

	// Set options
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}
	if o.logger != nil {
		db = &loggerDB{DB: db, logger: o.logger}
	}

	return db, nil
}
