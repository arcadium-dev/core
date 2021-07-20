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
	"database/sql"
)

func Open(config Config, opts ...Option) (DB, error) {
	sqldb, err := sql.Open(config.DriverName(), config.DSN())
	if err != nil {
		return nil, err
	}
	var db DB = &sqlDB{DB: sqldb}

	// Set options
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}

	// If there is a migration, run it against the sql db.
	if o.migration != nil {
		if err := o.migration(sqldb); err != nil {
			return nil, err
		}
	}
	return db, nil
}
