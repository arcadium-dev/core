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
	"net/url"

	"github.com/kelseyhightower/envconfig"

	"arcadium.dev/core/errors"
)

type (
	// Postgres holds the configuration settings needed to connect to a postgres database.
	// The DB, HOST and SSLMODE variables are required. The SSLMODE defaults to "enabled".
	//
	// For sslmode setting, see https://www.postgresql.org/docs/current/libpq-ssl.html
	Postgres struct {
		db             string // <PREFIX_>POSTGRES_DB
		user           string // <PREFIX_>POSTGRES_USER
		password       string // <PREFIX_>POSTGRES_PASSWORD
		host           string // <PREFIX_>POSTGRES_HOST
		port           string // <PREFIX_>POSTGRES_PORT
		sslMode        string // <PREFIX_>POSTGRES_SSLMODE
		sslCert        string // <PREFIX_>POSTGRES_SSLCERT
		sslKey         string // <PREFIX_>POSTGRES_SSLKEY
		sslRootCert    string // <PREFIX_>POSTGRES_SSLROOTCERT
		connectTimeout string // <PREFIX_>POSTGRES_CONNECT_TIMEOUT
	}
)

// NewPostgres returns the postgres configuration.
func NewPostgres() (*Postgres, error) {
	config := struct {
		DB             string `required:"true"`
		User           string
		Password       string
		Host           string `required:"true"`
		Port           string
		SSLMode        string `default:"verify-full"`
		SSLCert        string
		SSLKey         string
		SSLRootCert    string
		Level          string
		File           string
		Format         string
		ConnectTimeout string `split_words:"true"`
	}{}
	if err := envconfig.Process("postgres", &config); err != nil {
		return nil, errors.Wrap(err, "failed to load postgres configuration")
	}
	return &Postgres{
		db:             config.DB,
		user:           config.User,
		password:       config.Password,
		host:           config.Host,
		port:           config.Port,
		sslMode:        config.SSLMode,
		sslCert:        config.SSLCert,
		sslKey:         config.SSLKey,
		sslRootCert:    config.SSLRootCert,
		connectTimeout: config, Timeout,
	}, nil
}

// DSN returns a connection string corresponding to the postgres configuration.
//
// See https://godoc.org/github.com/lib/pq for connection string parameters.
func (p *Postgres) DSN() string {
	// Build the url
	u := &url.URL{Scheme: "postgres"}
	if p.DB != "" {
		u.Path = p.DB
	}
	if p.User != "" {
		u.User = url.UserPassword(p.User, p.Password)
	}
	host := ""
	if p.Host != "" {
		host = p.Host
	}
	if p.Port != "" {
		host += ":" + p.Port
	}
	u.Host = host

	// Build the query
	q := u.Query()
	if p.ConnectTimeout != "" {
		q.Add("connect_timeout", p.ConnectTimeout)
	}
	if p.SSLMode != "" {
		q.Add("sslmode", p.SSLMode)
	}
	if p.SSLCert != "" {
		q.Add("sslcert", p.SSLCert)
	}
	if p.SSLKey != "" {
		q.Add("sslkey", p.SSLKey)
	}
	if p.SSLRootCert != "" {
		q.Add("sslrootcert", p.SSLRootCert)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
