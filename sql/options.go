// Copyright 2021-2023 arcadium.dev <info@arcadium.dev>
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

// Package sql provides a database/sql wrapper.
package sql // import "arcadium.dev/core/sql"

type (
	// Option provides options for the creation of a db connection.
	Option interface {
		Apply(*Options)
	}

	// Options holds the db options.
	Options struct {
		ReconnectEnabled bool
		ReconnectRetries int
		TxRetries        int
	}
)

// WithReconnect will configure the db connection to reconnect if the
// connection has been closed due to an admin shutdown.
func WithReconnect(numRetries int) Option {
	return newOption(func(o *Options) {
		o.ReconnectEnabled = true
		if numRetries > 0 {
			o.ReconnectRetries = numRetries
		}
	})
}

// WithTxRetries will configure the db to retry failed transactions.
func WithTxRetries(numRetries int) Option {
	return newOption(func(o *Options) {
		if numRetries > 0 {
			o.TxRetries = numRetries
		}
	})
}

type (
	// option implements the Option interface.
	option struct {
		f func(*Options)
	}
)

func newOption(f func(*Options)) *option {
	return &option{f: f}
}

func (o *option) Apply(opts *Options) {
	o.f(opts)
}
