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

// Package middleware provides a set of middleware for the http server.
package middleware // import "arcadium.dev/core/http/middleware"

import (
	"net/http"

	"github.com/rs/zerolog"
)

type (
	Logging struct {
		Logger *zerolog.Logger
	}
)

// Requests is middleware to create a request specific logger, which includes
// the request's method and url as fields. The logger is passed through the
// request's context This also logs the incoming request.
func (l Logging) Requests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := l.Logger.With().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Logger()

		req := r.Clone(logger.WithContext(r.Context()))

		logger.Debug().Msg("request start")
		next.ServeHTTP(w, req)
		logger.Debug().Msg("request complete")
	})
}
