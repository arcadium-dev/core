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
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type (
	Postgres struct {
		DB             string // POSTGRES_DB
		User           string // POSTGRES_USER
		Password       string // POSTGRES_PASSWORD
		Host           string // POSTGRES_HOST
		Port           string // POSTGRES_PORT
		ConnectTimeout string `split_words:"true"` // POSTGRES_CONNECT_TIMEOUT
		SSLMode        string // POSTGRES_SSLMODE
		SSLCert        string // POSTGRES_SSLCERT
		SSLKey         string // POSTGRES_SSLKEY
		SSLRootCert    string // POSTGRES_SSLROOTCERT
	}
)

func NewPostgres() (*Postgres, error) {
	var p Postgres
	if err := envconfig.Process("postgres", &p); err != nil {
		return nil, errors.Wrap(err, "failed to load postgres configuration")
	}
	return &p, nil
}

// DSN returns a connection string corresponding to the postgres configuration.
//
// See https://godoc.org/github.com/lib/pq for connection string parameters.
func (p *Postgres) DSN() string {
	dsn := ""
	if p.DB != "" {
		dsn += fmt.Sprintf("dbname='%s' ", p.DB)
	}
	if p.User != "" {
		dsn += fmt.Sprintf("user='%s' ", p.User)
	}
	if p.Password != "" {
		dsn += fmt.Sprintf("password='%s' ", p.Password)
	}
	if p.Host != "" {
		dsn += fmt.Sprintf("host='%s' ", p.Host)
	}
	if p.Port != "" {
		dsn += fmt.Sprintf("port='%s' ", p.Port)
	}
	if p.ConnectTimeout != "" {
		dsn += fmt.Sprintf("connect_timeout='%s' ", p.ConnectTimeout)
	}
	if p.SSLMode != "" {
		dsn += fmt.Sprintf("sslmode='%s' ", p.SSLMode)
	}
	if p.SSLCert != "" {
		dsn += fmt.Sprintf("sslcert='%s' ", p.SSLCert)
	}
	if p.SSLKey != "" {
		dsn += fmt.Sprintf("sslkey='%s' ", p.SSLKey)
	}
	if p.SSLRootCert != "" {
		dsn += fmt.Sprintf("sslrootcert='%s' ", p.SSLRootCert)
	}
	return strings.TrimSpace(dsn)
}
