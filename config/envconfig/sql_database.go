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

package envconfig // import "arcadium.dev/core/config/envconfig

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	sqlPrefix = "sql_database"
)

type (
	DataSourceNamer interface {
		DSN() string
	}

	SQLDatabase struct {
		driver string // <PREFIX_>SQL_DATABASE_DRIVER
		DataSourceNamer
	}
)

func NewSQLDatabase(opts ...Option) (*SQLDatabase, error) {
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}
	prefix := o.prefix + sqlPrefix

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

func (db *SQLDatabase) DriverName() string {
	return db.driver
}
