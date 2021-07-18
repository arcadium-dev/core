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

package envconfig // import "arcadium.dev/core/config/envconfig

import (
	"strings"

	"arcadium.dev/core/config"
	"github.com/kelseyhightower/envconfig"

	"arcadium.dev/core/errors"
)

const (
	sqlPrefix = "sql_database"
)

type (
	// DataSourceNamer is the interface that wraps the DSN method.
	DataSourceNamer interface {
		DSN() string
	}

	// SQLDatabase holds the configuration information for an SQL database.
	SQLDatabase struct {
		driver string // <PREFIX_>SQL_DATABASE_DRIVER
		DataSourceNamer
	}
)

// NewSQLDatabase returns the configuration of an SQL database.
func NewSQLDatabase(opts ...config.Option) (*SQLDatabase, error) {
	o := &config.Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	prefix := o.Prefix() + sqlPrefix

	config := struct {
		Driver string `default:"postgres"`
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return nil, errors.Wrapf(err, "failed to load %s configuration", prefix)
	}

	var namer DataSourceNamer
	switch strings.TrimSpace(strings.ToLower(config.Driver)) {
	case "postgres":
		p, err := NewPostgres()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load postgres configuration")
		}
		namer = p
	case "mysql":
		fallthrough
	case "sqlite":
		fallthrough
	default:
		return nil, errors.Errorf("unsupported database driver: %s", config.Driver)
	}

	return &SQLDatabase{
		driver:          config.Driver,
		DataSourceNamer: namer,
	}, nil
}

// DriverName returns the name of the sql database driver.
func (db *SQLDatabase) DriverName() string {
	return db.driver
}
