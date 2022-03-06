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
	"net/url"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type (
	// SQL holds the configuration settings needed to connect to an sql database.
	SQL struct {
		url string
	}
)

const (
	sqlPrefix = "sql"
)

// NewSQL returns the sql configuration.
func NewSQL(opts ...Option) (SQL, error) {
	o := &Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	prefix := o.Prefix + sqlPrefix

	config := struct {
		URL string `required:"true"`
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return SQL{}, fmt.Errorf("failed to load %s configuration: %w", prefix, err)
	}
	if _, err := url.Parse(config.URL); err != nil {
		return SQL{}, fmt.Errorf("failed to parse %s connection url: %w", prefix, err)
	}
	return SQL{
		url: strings.TrimSpace(config.URL),
	}, nil
}

// URL returns the connection URL for the SQL database. The value is set from
// the <PREFIX_>SQL_URL environment variable.
func (s SQL) URL() string {
	return s.url
}
