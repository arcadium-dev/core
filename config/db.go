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

package config // import "arcadium.dev/core/config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type (
	// DB holds the configuration settings needed to connect to a database.
	DB struct {
		driver string
		dsn    string
	}
)

const (
	dbPrefix = "db"
)

// NewDB returns the db configuration.
func NewDB(opts ...Option) (DB, error) {
	o := &Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	prefix := o.Prefix + dbPrefix

	config := struct {
		Driver string `required:"true"`
		DSN    string `required:"true"`
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return DB{}, fmt.Errorf("failed to load %s configuration: %w", prefix, err)
	}
	return DB{
		driver: strings.TrimSpace(config.Driver),
		dsn:    strings.TrimSpace(config.DSN),
	}, nil
}

// Driver returns the database driver. The value is set from the
// <PREFIX_>DB_DRIVER environment variable.
func (db DB) Driver() string {
	return db.driver
}

// DSN returns the datasource name for the database. The value is set from
// the <PREFIX_>DB_DSN environment variable.
func (db DB) DSN() string {
	return db.dsn
}
