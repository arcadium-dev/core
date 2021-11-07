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

package config // import "arcadium.dev/core/config/config"

import (
	"strings"

	"github.com/kelseyhightower/envconfig"

	"arcadium.dev/core/errors"
)

const (
	sqlPrefix = "sql_database"
)

type (
	// DataSourceNamer is the interface that wraps the DSN method.
	DataSourceNamer interface {
		// DSN returns the data source name of the sql database.
		DSN() string
	}

	// SQLDatabase holds the configuration information for an SQL database.
	SQLDatabase interface {
		// DriverName returns the name of the sql database driver.
		DriverName() string
		DataSourceNamer
	}
)

// NewSQLDatabase returns the configuration of an SQL database.
func NewSQLDatabase(opts ...Option) (SQLDatabase, error) {
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}
	prefix := o.prefix + sqlPrefix

	config := struct {
		Driver string `default:"pgx"`
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return nil, errors.Wrapf(err, "failed to load %s configuration", prefix)
	}

	var namer DataSourceNamer
	driver := strings.TrimSpace(strings.ToLower(config.Driver))
	switch driver {
	case "pgx", "postgres":
		p, err := NewPostgres(opts...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load postgres configuration")
		}
		namer = p
	case "mysql":
		fallthrough
	case "sqlite":
		fallthrough
	default:
		return nil, errors.Errorf("unsupported database driver: %s", driver)
	}

	return &sqlDatabase{
		driver:          driver,
		DataSourceNamer: namer,
	}, nil
}

type (
	sqlDatabase struct {
		driver string // <PREFIX_>SQL_DATABASE_DRIVER
		DataSourceNamer
	}
)

func (db *sqlDatabase) DriverName() string {
	return db.driver
}
